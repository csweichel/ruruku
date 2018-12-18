import { AppStateContent, getClient } from 'src/types/app-state';
import { TestRunStatus } from 'src/api/v1/session_pb';
import * as React from 'react';
import { List, Header, Label } from 'semantic-ui-react';

export interface SessionDetailViewProps {
    appState: AppStateContent;
}

interface SessionDetailViewState {
    session?: TestRunStatus;
}

export class SessionDetailView extends React.Component<SessionDetailViewProps, SessionDetailViewState> {

    constructor(props: SessionDetailViewProps) {
        super(props);
        this.state = {};
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
                        <Header>Annotations:</Header>
                        {content}
                    </List.Content>
                </List.Item>
            );
        }

        return (
            <List celled={true}>
                <List.Item>
                    <List.Content>
                        <Header>Name:</Header> {session.getName()}
                    </List.Content>
                </List.Item>
                {annotations}
            </List>
        );
    }

    protected async fetchState() {
        const session = await getClient(this.props.appState).getStatus();
        this.setState({ session });
    }

}