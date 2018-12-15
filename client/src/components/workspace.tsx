import * as React from 'react';
import logo from '../logo.svg';
// import { TestplanView } from './testplan-view';
import { Sidebar, Segment, Message } from 'semantic-ui-react';
// import { Participant } from '../types/participant';
// import { TestRunStatus, SessionStatusRequest, SessionStatusResponse, Testcase, TestRunState, ClaimRequest, SessionUpdatesRequest, SessionUpdateResponse, ContributionRequest } from '../api/v1/session_pb';
// import { grpc } from 'grpc-web-client';
// import { SessionService } from '../api/v1/session_pb_service';
// import { HOST } from '../api/host';
// import { TestplanView } from './testplan-view';
import { AppStateContent } from 'src/types/app-state';
import { TestplanView } from './testplan-view';


export interface WorkspaceProps {
    appState: AppStateContent
}

interface WorkspaceState {
    sidebar?: any
    error?: string
}

export class Workspace extends React.Component<WorkspaceProps, WorkspaceState> {
    constructor(props: WorkspaceProps) {
        super(props);

        this.showSidebar = this.showSidebar.bind(this);

        this.state = { };
    }

    public render() {
        const error = this.state.error
            ? <Message error={true}>{this.state.error}</Message>
            : undefined;

        return (
            <div className="workspace">
                <div id="header">
                    <img src={logo} className="app-logo" alt="logo" />
                    <div className="info">
                        <div className="session">{this.props.appState.session}</div>
                        <div className="username" />
                    </div>
                </div>

                <div className="main">
                    {error}
                    <Sidebar.Pushable as={Segment} attached="bottom" className="no-border">
                        <Sidebar width="very wide" animation="overlay" visible={!!this.state.sidebar} icon="labeled" vertical={true} inline={true} inverted={false} direction="right">
                            {this.state.sidebar}
                        </Sidebar>
                        <Sidebar.Pusher>
                            <Segment basic={true} className="no-padding">
                                <TestplanView
                                    appState={this.props.appState}
                                    showDetails={this.showSidebar} />
                            </Segment>
                        </Sidebar.Pusher>
                    </Sidebar.Pushable>
                </div>
            </div>
        )
    }

    protected showSidebar(sidebar: any | undefined) {
        this.setState({ sidebar });
    }

}