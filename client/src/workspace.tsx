import * as React from 'react';
import { WelcomeRequest, TestSuite, TestRun, TestParticipant } from '../../protocol/protocol'
import logo from './logo.svg';
import { TestplanView } from './testplan-view';
import { Button } from 'semantic-ui-react';

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

    constructor(props: WorkspaceProps) {
        super(props);
        props.ws.onmessage = this.onMessage.bind(this);

        this.sendWelcome();
    }

    public componentWillMount() {
        this.setState({ view: "none" });
        this.sendWelcome();
    }

    public render() {
        let main = <div />;
        if(this.state.view === "plan") {
            main = <TestplanView suite={this.state.suite} run={this.state.run} participant={this.state.participant} />;
        }

        return (
            <div className="workspace">
                <div className="header">
                    <img src={logo} className="app-logo" alt="logo" />
                </div>
                <div className="sidebar">
                    <Button label="Test Plan" />
                    <ul>
                        <li><header><a>Plan</a></header></li>
                        <li>
                            <header>Your tests</header>
                            <ul>
                                <li><a>Item 1</a></li>
                            </ul>
                        </li>
                    </ul>
                </div>
                <div className="main">
                    {main}
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
        }
    }

    protected sendWelcome() {
        const welcome: WelcomeRequest = {
            type: "welcome",
            name: this.props.name
        };
        this.props.ws.send(JSON.stringify(welcome));
    }

}