import * as React from 'react';
import logo from '../logo.svg';
import { Sidebar, Segment, Message, TransitionGroup } from 'semantic-ui-react';
import { AppStateContent } from 'src/types/app-state';
import { TestplanView } from './testplan-view';

import './workspace.css';

export interface WorkspaceProps {
    appState: AppStateContent
}

interface WorkspaceState {
    sidebar?: any
}

export class Workspace extends React.Component<WorkspaceProps, WorkspaceState> {
    constructor(props: WorkspaceProps) {
        super(props);

        this.showSidebar = this.showSidebar.bind(this);

        this.state = { };
    }

    public render() {
        const error = this.props.appState.error
            ? <Message error={true}>{this.props.appState.error}</Message>
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
                    <TransitionGroup animation="drop">
                        {error}
                    </TransitionGroup>
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