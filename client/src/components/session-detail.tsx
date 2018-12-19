import { AppStateContent, getClient } from 'src/types/app-state';
import { TestRunStatus } from 'src/api/v1/session_pb';
import * as React from 'react';
import { List, Header, Label, Segment, Grid, Transition } from 'semantic-ui-react';

export interface SessionDetailViewProps {
    appState: AppStateContent;
    onChangePwd: () => void;
}

interface SessionDetailViewState {
    session?: TestRunStatus;
}

export class SessionDetailView extends React.Component<SessionDetailViewProps, SessionDetailViewState> {

    constructor(props: SessionDetailViewProps) {
        super(props);
        this.state = {
        };
    }

    public componentWillMount() {
        this.fetchState();
    }

    public render() {
        const { session } = this.state;
        if (!session) {
            return <div />;
        }

        let annotations: JSX.Element | undefined;
        if (session.getAnnotationsMap().getLength() > 0) {
            const content: JSX.Element[] = [];
            const annotationMap = session.getAnnotationsMap();
            for (const kv of annotationMap.toArray()) {
                const name = kv[0];
                const value = kv[1];
                content.push(<Label key={name}>{name}<Label.Detail>{value}</Label.Detail></Label>);
            }

            annotations = (
                <List.Item>
                    <List.Content>
                        <Header subheader={true}>Annotations</Header>
                        {content}
                    </List.Content>
                </List.Item>
            );
        }

        return (
            <Segment>
                <Transition.Group animation="fly left">
                    <Grid columns={2} relaxed='very'>
                <Grid.Column>
                    <List>
                        <List.Item>
                            <List.Content>
                                <Header subheader={true}>User</Header> {(this.props.appState.user || { name: "" }).name}
                            </List.Content>
                        </List.Item>
                        <List.Item>
                            <List.Content>
                                <Header subheader={true}>Session Name</Header>
                                {session.getName()} / {session.getId()}
                            </List.Content>
                        </List.Item>
                        {annotations}
                    </List>
                </Grid.Column>
                <Grid.Column>
                    <List selection={true} divided={true}>
                        <List.Item onClick={this.props.onChangePwd}>
                            <List.Content>Change Password</List.Content>
                        </List.Item>
                        <List.Item onClick={this.props.appState.resetSession}>
                            <List.Content>Switch Session</List.Content>
                        </List.Item>
                        <List.Item onClick={this.props.appState.logout}>
                            <List.Content>Logout</List.Content>
                        </List.Item>
                    </List>
                </Grid.Column>
            </Grid>
                </Transition.Group>
            </Segment>
        );
    }

    protected async fetchState() {
        const session = await getClient(this.props.appState).getStatus();
        this.setState({ session });
    }

}