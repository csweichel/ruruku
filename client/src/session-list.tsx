import * as React from 'react';
import { Message, Dropdown, Form, Segment, Input, Button, InputProps, DropdownProps } from 'semantic-ui-react';
import { grpc } from 'grpc-web-client';
import { SessionService } from './api/v1/api_pb_service';
import { ListSessionsRequest, ListSessionsResponse, RegistrationRequest, RegistrationResponse } from './api/v1/api_pb';
import { HOST } from './api/host';
import "./session-list.css";
import { Participant } from './types/participant';

export interface SessionListProps {
    onJoin: (participant: Participant) => void
}

interface SessionListState {
    error?: string
    sessions?: SessionDescription[]
    selectedSessionID?: string
    participantName?: string
}

interface SessionDescription {
    id: string
    name: string
    isOpen: boolean
}

export class SessionList extends React.Component<SessionListProps, SessionListState> {
    protected poller?: NodeJS.Timeout;

    constructor(props: SessionListProps) {
        super(props);
        this.state = {}

        this.setParticipantName = this.setParticipantName.bind(this);
        this.selectSession = this.selectSession.bind(this);
        this.cannotJoin = this.cannotJoin.bind(this);
        this.joinSession = this.joinSession.bind(this);
    }

    public componentWillMount() {
        this.getSessions();
        this.poller = setInterval(this.getSessions.bind(this), 5000);
    }

    public componentWillUnmount() {
        if (this.poller) {
            clearInterval(this.poller);
            this.poller = undefined;
        }
    }

    public render() {
        let body: JSX.Element;
        if (this.state.error) {
            body = <Message error={true}>{this.state.error}</Message>;
        } else if (this.state.sessions) {
            const options = this.state.sessions.filter(s => s.isOpen).map(s => { return {
                value: s.id,
                text: s.name
            }})

            body = (
                <Segment>
                    <Form>
                        <Form.Field>
                            <Dropdown borderless={false} fluid={true} options={options} placeholder="Session" onChange={this.selectSession} />
                        </Form.Field>
                        <Form.Field>
                            <Input type="text" placeholder="Your name" value={this.state.participantName} onChange={this.setParticipantName} />
                        </Form.Field>
                        <Form.Field>
                            <Button type="submit" primary={true} disabled={this.cannotJoin()} onClick={this.joinSession}>Join</Button>
                        </Form.Field>
                    </Form>
                </Segment>
            )
        } else {
            body = <div>Connecting ...</div>;
        }

        return <div id="selectSession">{body}</div>
    }

    protected cannotJoin() {
        return !this.state.participantName || !this.state.selectedSessionID;
    }

    protected selectSession(evt: React.SyntheticEvent, props: DropdownProps) {
        this.setState({ selectedSessionID: props.value as string });
    }

    protected setParticipantName(ev: React.SyntheticEvent, props: InputProps) {
        this.setState({ participantName: props.value as string })
    }

    protected joinSession(evt: React.SyntheticEvent) {
        evt.preventDefault();

        try {
            const req = new RegistrationRequest();
            req.setSessionid(this.state.selectedSessionID!);
            req.setName(this.state.participantName!);

            grpc.invoke(SessionService.Register, {
                request: req,
                host: HOST,
                onMessage: msg => {
                    const resp = msg as RegistrationResponse;
                    this.props.onJoin({
                        sessionID: this.state.selectedSessionID!,
                        name: this.state.participantName!,
                        token: resp.getToken()
                    });
                },
                onEnd: res => {
                    // nothing to do here
                },
            });
        } catch (err) {
            console.error("Error while joining session", err);
            this.setState({ error: "Error: " + err });
        }
    }

    protected getSessions() {
        const sessions: SessionDescription[] = [];
        try {
            grpc.invoke(SessionService.List, {
            request: new ListSessionsRequest(),
                host: HOST,
                onMessage: msg => {
                    const resp = msg as ListSessionsResponse;
                    sessions.push({
                        id: resp.getId(),
                        name: resp.getName(),
                        isOpen: resp.getIsopen()
                    });
                },
                onEnd: res => {
                    this.setState({sessions})
                }
            });
        } catch (err) {
            console.error("Error while getting session list", err);
            this.setState({ error: err.toString() });
        }
    }

}