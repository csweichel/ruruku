import * as React from 'react';
import { WelcomeRequest, TestSuite, TestRun, TestParticipant, TestCase, ClaimRequest } from '../../protocol/protocol'
import logo from './logo.svg';
import { TestplanView } from './testplan-view';

export interface WorkspaceProps {
    ws: WebSocket
    name: string
}

interface WorkspaceState {
    suite: TestSuite
    run: TestRun
    participant: TestParticipant
    view: "none" | "plan"
}

export class Workspace extends React.Component<WorkspaceProps, WorkspaceState> {
    protected keepAliveDisposable?: () => void;

    constructor(props: WorkspaceProps) {
        super(props);

        // ES6 classes do not autobind to this
        props.ws.onmessage = this.onMessage.bind(this);
        this.claimTestCase = this.claimTestCase.bind(this);

        this.sendWelcome();
        this.keepAliveDisposable = this.keepConnectionAlive();
    }

    public componentWillMount() {
        this.setState({ view: "none" });
        this.sendWelcome();
    }

    public componentWillUnmount() {
        if (this.keepAliveDisposable) {
            this.keepAliveDisposable();
        }
    }

    public render() {
        return (
            <div className="workspace">
                <div className="header">
                    <img src={logo} className="app-logo" alt="logo" />
                    <div className="username">{this.props.name}</div>
                </div>

                <div className="main">
                    <TestplanView
                        suite={this.state.suite}
                        run={this.state.run}
                        participant={this.state.participant}
                        claimTestCase={this.claimTestCase} />
                </div>
            </div>
        )
    }

    protected onMessage(ev: MessageEvent) {
        const msg = JSON.parse(ev.data);
        if (!msg) {
            console.warn("Received unparseable message", ev.data);
        }

        if (msg.type === "welcome") {
            this.setState({
                view: "plan",
                suite: msg.suite,
                run: msg.run,
                participant: msg.participant
            });
        } else if (msg.type === "update") {
            this.setState({
                run: msg.run,
                participant: msg.participant
            });
        }
    }

    protected sendWelcome() {
        const welcome: WelcomeRequest = {
            type: "welcome",
            name: this.props.name
        };
        this.props.ws.send(JSON.stringify(welcome));
    }

    protected claimTestCase(testCase: TestCase, claim: boolean) {
        const msg: ClaimRequest = {
            type: "claim",
            caseId: `${testCase.group}/${testCase.id}`,
            claim
        };
        this.props.ws.send(JSON.stringify(msg));
    }

    protected keepConnectionAlive(): () => void {
        const timeout = setInterval(() => {
            try {
                this.props.ws.send(JSON.stringify({ type: "keep-alive" }));
            } catch(err) {
                // ignore for now until we have better connection management
            }
        }, 10000);
        return () => clearTimeout(timeout);
    }

}