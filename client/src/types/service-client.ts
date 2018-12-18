import { AuthenticationRequest, AuthenticationRespose, RenewTokenRequest, RenewTokenResponse } from '../api/v1/user_pb';
import { grpc } from 'grpc-web-client';
import { UserService } from '../api/v1/user_pb_service';
import { ProtobufMessage } from 'grpc-web-client/dist/typings/message';
import { MethodDefinition } from 'grpc-web-client/dist/typings/service';
import { ClaimRequest, TestRunState, ContributionRequest, TestRunStatus, SessionStatusRequest, SessionStatusResponse, SessionUpdatesRequest, SessionUpdateResponse, RegistrationRequest } from '../api/v1/session_pb';
import { SessionService } from '../api/v1/session_pb_service';
import { AppStateContent, getAuthentication } from './app-state';

export interface InvokeOptions<Response, TResponse> {
    onMessage?: (response: TResponse, memento?: Response) => Response;
    onEnd?: (code: grpc.Code, resolve: (res: Response) => void, reject: (err: Error) => void, memento?: Response) => void
    isOK?: (code: grpc.Code) => boolean;
}

export interface Disposable {
    dispose(): void
}

export class Client {

    public constructor(protected readonly appState: AppStateContent, protected readonly host: string) {
    }

    public async login(username: string, password: string): Promise<string> {
        const req = new AuthenticationRequest();
        req.setUsername(username);
        req.setPassword(password);

        const result = await this.invoke(UserService.AuthenticateCredentials, req, {
            onMessage: (msg: AuthenticationRespose) => {
                return msg.getToken();
            }
        });
        this.appState.user = { name: username };
        this.appState.token = result;
        return result!;
    }

    public async joinSession(session: string): Promise<void> {
        const req = new RegistrationRequest();
        req.setSession(session);
        await this.invoke(SessionService.Register, req);
    }

    public async claim(testcase: string, claim: boolean): Promise<void> {
        const req = new ClaimRequest();
        req.setSession(this.appState.session!);
        req.setTestcaseid(testcase);
        req.setClaim(claim);
        await this.invoke(SessionService.Claim, req);
    }

    public async contribute(testcase: string, result: TestRunState, comment: string): Promise<void> {
        const req = new ContributionRequest();
        req.setSession(this.appState.session!);
        req.setTestcaseid(testcase);
        req.setResult(result);
        req.setComment(comment);
        await this.invoke(SessionService.Contribute, req);
    }

    public async renewToken(): Promise<void> {
        const resp = await this.invoke(UserService.RenewToken, new RenewTokenRequest(), {
            onMessage: (msg: RenewTokenResponse) => msg.getToken()
        });
        this.appState.token = resp;
    }

    public async getStatus(): Promise<TestRunStatus> {
        const req = new SessionStatusRequest();
        req.setId(this.appState.session!);
        const result = await this.invoke(SessionService.Status, req, {
            onMessage: (msg: SessionStatusResponse) => msg.getStatus()
        });
        return result!;
    }

    public async listenForUpdates(callback: UpdateListenerCallback): Promise<Disposable> {
        const listener = new UpdateListener(this.host, getAuthentication(this.appState)!, callback);
        await listener.start(this.appState.session!);
        return listener;
    }

    protected async invoke<Response, TRequest extends ProtobufMessage, TResponse extends ProtobufMessage, M extends MethodDefinition<TRequest, TResponse>>(methodDescriptor: M, req: TRequest, options?: InvokeOptions<Response, TResponse>): Promise<Response | undefined> {
        return new Promise<Response>((resolve, reject) => {
            const props = options || {};
            try {
                let memento: Response | undefined;
                grpc.invoke(methodDescriptor, {
                    request: req,
                    host: this.host,
                    metadata: getAuthentication(this.appState),
                    onMessage: (msg: TResponse) => {
                        if (props.onMessage) {
                            memento = props.onMessage(msg, memento);
                        }
                    },
                    onEnd: (res, msg) => {
                        const isOK = props.isOK || ((c: grpc.Code) => c === grpc.Code.OK);
                        if (!isOK(res)) {
                            console.warn("Request failed", res, msg);
                            reject(msg);
                            return;
                        }

                        if (props.onEnd) {
                            props.onEnd(res, resolve, reject, memento);
                        } else {
                            resolve(memento);
                        }
                    }
                });
            } catch(err) {
                console.warn("Request failed", err);
                reject(err);
            }
        });
    }

}

export type UpdateListenerCallback = (status?: TestRunStatus, err?: Error) => void;

class UpdateListener implements Disposable {
    protected dontReconnect: boolean;
    protected didEverConnect: boolean;
    protected client?: { close: () => void };

    constructor(
        protected readonly host: string,
        protected readonly metadata: grpc.Metadata,
        protected readonly callback: UpdateListenerCallback) {
    }

    public async start(session: string): Promise<void> {
        const req = new SessionUpdatesRequest();
        req.setId(session);

        return new Promise<void>((resolve, reject) => {
            try {
                this.client = grpc.invoke(SessionService.Updates, {
                    host: this.host,
                    metadata: this.metadata,
                    request: req,
                    onMessage: (msg: SessionUpdateResponse) => {
                        this.didEverConnect = true;
                        this.callback(msg.getStatus());
                        resolve();
                    },
                    onEnd: (res, msg) => {
                        if (res === grpc.Code.OK) {
                            return;
                        } else {
                            this.callback(undefined, new Error(msg));
                            this.reconnect(session);
                        }
                    }
                });
            } catch(err) {
                this.callback(undefined, err);
                if (this.didEverConnect) {
                    setTimeout(this.reconnect.bind(this, session), 500);
                } else {
                    reject(err);
                }
            }
        })
    }

    public dispose() {
        this.dontReconnect = true;
        if (this.client) {
            this.client.close();
        }
    }

    protected async reconnect(session: string) {
        if (this.dontReconnect) {
            return;
        }

        this.start(session);
    }

}