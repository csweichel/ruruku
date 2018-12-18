import { AppStateContent, getAuthentication, getClient } from '../types/app-state';
import * as React from 'react';
import { Segment, Form, Dropdown, Button, DropdownProps, Message } from 'semantic-ui-react';
import { grpc } from 'grpc-web-client';
import { SessionService } from '../api/v1/session_pb_service';
import { ListSessionsRequest, ListSessionsResponse } from '../api/v1/session_pb';
import { HOST } from '../api/host';

export interface SessionSelectorProps {
    appState: AppStateContent
    onSelect: (session: string) => void
}

interface SessionSelectorState {
    sessions?: SessionDescription[]
    selectedSessionID?: string
}

interface SessionDescription {
    id: string
    name: string
    isOpen: boolean
}

export class SessionSelector extends React.Component<SessionSelectorProps, SessionSelectorState> {
    protected poller?: NodeJS.Timeout;

    constructor(props: SessionSelectorProps) {
        super(props);

        this.state = {};

        this.selectSession = this.selectSession.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
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
        const options = (this.state.sessions || []).filter(s => s.isOpen).map(s => { return {
            value: s.id,
            text: s.name
        }});
        if (options.length === 1) {
            this.props.onSelect(options[0].value);
        }

        let error: JSX.Element | undefined;
        if (this.props.appState.error) {
            error = <Message error={true}>{this.props.appState.error}</Message>
        };

        return (
            <Segment id="login-form">
                {error}
                <Form onSubmit={this.onSubmit}>
                    <Form.Field>
                        <Dropdown selection={true} fluid={true} options={options} placeholder="Session" onChange={this.selectSession} />
                    </Form.Field>
                    <Form.Field>
                        <Button type="submit" primary={true} disabled={!this.state.selectedSessionID}>Join</Button>
                    </Form.Field>
                </Form>
            </Segment>
        );
    }

    protected onSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        getClient(this.props.appState).joinSession(this.state.selectedSessionID!);
        this.props.onSelect(this.state.selectedSessionID!);
    }

    protected selectSession(evt: React.SyntheticEvent, props: DropdownProps) {
        this.setState({ selectedSessionID: props.value as string });
    }

    protected getSessions() {
        const sessions: SessionDescription[] = [];
        try {
            grpc.invoke(SessionService.List, {
            request: new ListSessionsRequest(),
                host: HOST,
                metadata: getAuthentication(this.props.appState),
                onMessage: (resp: ListSessionsResponse) => {
                    sessions.push({
                        id: resp.getId(),
                        name: resp.getName(),
                        isOpen: resp.getIsopen()
                    });
                },
                onEnd: (res, msg) => {
                    this.setState({sessions});
                    if (res === grpc.Code.Unknown) {
                        // really gotta figure out where these unknown EOT events come from
                    } else if (res !== grpc.Code.OK) {
                        this.props.appState.setError(msg);
                    }
                }
            });
        } catch (err) {
            console.error("Error while getting session list", err);
            this.props.appState.setError(err);
        }
    }

}