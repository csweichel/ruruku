import * as React from 'react';
import { WelcomeRequest, TestSuite, TestRun, TestParticipant, TestCase, ClaimRequest, TestCaseResult, NewTestCaseRunRequest } from '../../protocol/protocol'
import logo from './logo.svg';
import { TestplanView } from './testplan-view';
import { Sidebar, Segment } from 'semantic-ui-react';

export interface WorkspaceProps {
    ws: WebSocket
    name: string
}

interface WorkspaceState {
    suite: TestSuite
    run: TestRun
    participant: TestParticipant
    view: "none" | "plan"
    sidebar: any | undefined
}

export class Workspace extends React.Component<WorkspaceProps, WorkspaceState> {
    protected keepAliveDisposable?: () => void;

    constructor(props: WorkspaceProps) {
        super(props);

        // ES6 classes do not autobind to this
        props.ws.onmessage = this.onMessage.bind(this);
        this.claimTestCase = this.claimTestCase.bind(this);
        this.submitTestcaseRun = this.submitTestcaseRun.bind(this);
        this.showSidebar = this.showSidebar.bind(this);

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
                <div id="header">
                    <img src={logo} className="app-logo" alt="logo" />
                    <div className="username">{this.props.name}</div>
                </div>

                <div className="main">
                    <Sidebar.Pushable as={Segment} attached="bottom" className="no-border">
                        <Sidebar animation="overlay" visible={!!this.state.sidebar} icon="labeled" vertical={true} inline={true} inverted={false} direction="right">
                            {this.state.sidebar}
                        </Sidebar>
                        <Sidebar.Pusher>
                            <Segment basic={true} className="no-padding">
                                <TestplanView
                                    suite={this.state.suite}
                                    run={this.state.run}
                                    participant={this.state.participant}
                                    claimTestCase={this.claimTestCase}
                                    submitTestCaseRun={this.submitTestcaseRun}
                                    showDetails={this.showSidebar} />
                            </Segment>
                        </Sidebar.Pusher>
                    </Sidebar.Pushable>
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

    protected submitTestcaseRun(testcase: TestCase, result: TestCaseResult, comment: string) {
        const msg: NewTestCaseRunRequest = {
            type: "newTestCaseRun",
            case: testcase.id,
            caseGroup: testcase.group,
            comment,
            result,
            start: new Date(),
        };
        this.props.ws.send(JSON.stringify(msg));
    }

    protected showSidebar(sidebar: any | undefined) {
        this.setState({ sidebar });
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