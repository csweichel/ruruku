import * as React from 'react';
import { WelcomeRequest } from '../../protocol/protocol'

export interface WorkspaceProps {
    ws: WebSocket
    name: string
}

interface WorkspaceState {
    foo?: boolean
}

export class Workspace extends React.Component<WorkspaceProps, WorkspaceState> {

    constructor(props: WorkspaceProps) {
        super(props);
        props.ws.onmessage = this.onMessage.bind(this);

        this.sendWelcome();
    }

    public onMessage(ws: WebSocket, ev: MessageEvent) {
        console.log(ev);
    }

    public render() {
        return (
            <div>{ this.props.name }</div>
        )
    }

    protected sendWelcome() {
        const welcome: WelcomeRequest = {
            type: "welcome",
            name: this.props.name
        };
        this.props.ws.send(JSON.stringify(welcome));
    }

}