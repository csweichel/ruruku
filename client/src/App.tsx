import * as React from 'react';
import './App.css';
import 'semantic-ui-css/semantic.min.css';

import logo from './logo.svg';
import { LoginForm } from './login-form';
import { Workspace } from './workspace';

type ConnectionState = "not-connected" | "connecting" | "connected";

interface AppState {
    name: string
    connection: ConnectionState
    socket?: WebSocket
    token: string
}

class App extends React.Component<{}, AppState> {

    constructor(props: {}) {
        super(props);

        this.connect = this.connect.bind(this);
    }

    public componentWillMount() {
        this.setState({
            name: "none",
            token: window.location.search.substring(1)
        });
    }

    public render() {
        let loginOrWorkspace;
        if (this.state.socket && this.state.connection === 'connected') {
            return <div className="app app-connected"><Workspace name={this.state.name} ws={this.state.socket} /></div>;
        } else if (this.state.connection === 'connecting') {
            loginOrWorkspace = <div>Connecting ...</div>
        } else if(this.state.token) {
            loginOrWorkspace = <LoginForm handleSubmit={this.connect} />
        } else {
            loginOrWorkspace = <b>No token present. Are you sure you have the right link?</b>
        }

        return (
            <div className={ this.state.socket ? "app app-connected" : "app app-login" }>
                <header id="header">
                    <img src={logo} className="app-logo" alt="logo" />
                </header>
                <div className="body">
                    { loginOrWorkspace }
                </div>
            </div>
        );
    }

    protected connect(name: string) {
        const token = this.state.token;
        if (!token) {
            throw new Error("No authentication token available. Cannot connect");
        }

        this.setState({ name, connection: 'connecting' });

        let protocol = 'ws';
        if (window.location.protocol === 'https:') {
            protocol = 'wss';
        }

        const ws = new WebSocket(`${protocol}://${window.location.host}/ws/${token}`);
        ws.onclose = (ev: CloseEvent) => {
            this.setState({ socket: undefined });
            if (this.state.connection === 'connected') {
                this.connect(name);
            } else {
                this.setState({ connection: 'not-connected' });
            }
        };
        ws.onopen = () => this.setState({ socket: ws, connection: 'connected' });
        ws.onerror = (err) => {
            console.log(err);

            this.setState({ socket: undefined });
            if (this.state.connection === 'connected') {
                this.connect(name);
            } else {
                alert("Cannot connect - make sure your link is correct");
                this.setState({ connection: 'not-connected' });
            }
        };
    }

}

export default App;
