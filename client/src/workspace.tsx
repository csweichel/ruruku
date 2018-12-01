import * as React from 'react';
import logo from './logo.svg';
// import { TestplanView } from './testplan-view';
import { Sidebar, Segment, Message } from 'semantic-ui-react';
import { Participant } from './types/participant';
import { TestRunStatus, SessionStatusRequest, SessionStatusResponse, Testcase, TestRunState, ClaimRequest } from './api/v1/api_pb';
import { grpc } from 'grpc-web-client';
import { SessionService } from './api/v1/api_pb_service';
import { HOST } from './api/host';
import { TestplanView } from './testplan-view';


export interface WorkspaceProps {
    participant: Participant
}

interface WorkspaceState {
    status?: TestRunStatus
    sidebar?: any
    error?: string
}

export class Workspace extends React.Component<WorkspaceProps, WorkspaceState> {
    protected keepAliveDisposable?: () => void;

    constructor(props: WorkspaceProps) {
        super(props);

        // ES6 classes do not autobind to this
        this.claimTestCase = this.claimTestCase.bind(this);
        this.submitTestcaseRun = this.submitTestcaseRun.bind(this);
        this.showSidebar = this.showSidebar.bind(this);
        this.getStatus = this.getStatus.bind(this);

        this.state = { };
    }

    public componentWillMount() {
        this.getStatus();
    }

    public render() {
        let session = "";
        let mainContent = <div>No testplan available</div>;
        if (this.state.status) {
            session = `${this.state.status.getName()}: ${this.state.status.getPlanid()}`;
            mainContent = <TestplanView
                status={this.state.status}
                participant={this.props.participant}
                claimTestCase={this.claimTestCase}
                submitTestCaseRun={this.submitTestcaseRun}
                showDetails={this.showSidebar} />;
        }

        const error = this.state.error
            ? <Message error={true}>{this.state.error}</Message>
            : undefined;

        return (
            <div className="workspace">
                <div id="header">
                    <img src={logo} className="app-logo" alt="logo" />
                    <div className="info">
                        <div className="session">{session}</div>
                        <div className="username">{this.props.participant.name}</div>
                    </div>
                </div>

                <div className="main">
                    {error}
                    <Sidebar.Pushable as={Segment} attached="bottom" className="no-border">
                        <Sidebar width="very wide" animation="overlay" visible={!!this.state.sidebar} icon="labeled" vertical={true} inline={true} inverted={false} direction="left">
                            {this.state.sidebar}
                        </Sidebar>
                        <Sidebar.Pusher>
                            <Segment basic={true} className="no-padding">
                                {mainContent}
                            </Segment>
                        </Sidebar.Pusher>
                    </Sidebar.Pushable>
                </div>
            </div>
        )
    }

    protected getStatus() {
        try {
            const req = new SessionStatusRequest();
            req.setId(this.props.participant.sessionID);

            grpc.invoke(SessionService.Status, {
                request: req,
                host: HOST,
                onMessage: msg => {
                    const resp = msg as SessionStatusResponse;
                    this.setState({ status: resp.getStatus() });
                },
                onEnd: res => {
                    // nothing to do here
                }
            });
        } catch (err) {
            console.error("Error while retrieving status", err);
            this.setState({ error: err.toString() });
        }
    }

    protected claimTestCase(testcaseId: string, participantToken: string, claim: boolean) {
        const req = new ClaimRequest();
        req.setParticipanttoken(participantToken);
        req.setTestcaseid(testcaseId);
        req.setClaim(claim);

        grpc.invoke(SessionService.Claim, {
            request: req,
            host: HOST,
            onEnd: res => {
                this.getStatus();
            }
        });
    }

    protected submitTestcaseRun(testCase: Testcase, participant: Participant, result: TestRunState, comment: string) {
        console.log("submitTestcaseRun", testCase, participant, result, comment);
    }

    protected showSidebar(sidebar: any | undefined) {
        this.setState({ sidebar });
    }

}