import * as React from 'react';
import logo from '../logo.svg';
import { Sidebar, Segment, Message, TransitionGroup, Accordion, Icon } from 'semantic-ui-react';
import { AppStateContent, getClient } from 'src/types/app-state';
import { TestplanView } from './testplan-view';

import './workspace.css';
import { TestRunStatus } from 'src/api/v1/session_pb';
import { SessionDetailView } from './session-detail';
import { ChangePasswordForm } from './change-password-form';

export interface WorkspaceProps {
    appState: AppStateContent
}

interface WorkspaceState {
    sidebar?: any
    cornerDropdown?: 'session-info' | 'chpwd';
    sessionInfo?: TestRunStatus;
}

export class Workspace extends React.Component<WorkspaceProps, WorkspaceState> {
    protected tokenRenewal?: NodeJS.Timer;

    constructor(props: WorkspaceProps) {
        super(props);

        this.showSidebar = this.showSidebar.bind(this);
        this.toggleCornerMenu = this.toggleCornerMenu.bind(this);

        this.state = { };
    }

    public componentWillMount() {
        this.startTokenRenewal();
    }

    public componentWillUnmount() {
        if (this.tokenRenewal) {
            clearTimeout(this.tokenRenewal);
        }
    }

    public render() {
        const error = this.props.appState.error
            ? <Message error={true}>{this.props.appState.error}</Message>
            : undefined;

        let content = <SessionDetailView appState={this.props.appState} onChangePwd={this.openCornerMenu.bind(this, 'chpwd')} />;
        if (this.state.cornerDropdown === 'chpwd') {
            content = <ChangePasswordForm appState={this.props.appState} onClose={this.openCornerMenu.bind(this, undefined)} />;
        }

        return (
            <div className="workspace">
                <div id="header">
                    <img src={logo} className="app-logo" alt="logo" />
                    <div className={"info " + (!!this.state.cornerDropdown ? "open" : "closed")}>
                        <Accordion styled={true}>
                            <Accordion.Title onClick={this.toggleCornerMenu}>
                                {!!this.state.cornerDropdown ? "Close" : (<Icon name="user circle" />)}
                            </Accordion.Title>
                            <Accordion.Content active={!!this.state.cornerDropdown}>
                                { content }
                            </Accordion.Content>
                        </Accordion>
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

    protected toggleCornerMenu() {
        this.setState({ cornerDropdown: !this.state.cornerDropdown ? 'session-info' : undefined });
    }

    protected openCornerMenu(cornerDropdown?: 'session-info' | 'chpwd') {
        this.setState({ cornerDropdown });
    }

    protected startTokenRenewal() {
        const renewer = async () => {
            try {
                await getClient(this.props.appState).renewToken();
            } catch (err) {
                this.props.appState.setError("Error while renewing user token: " + err);
            }
        };
        this.tokenRenewal = setInterval(renewer.bind(this), 5 * 60 * 1000);
    }

}