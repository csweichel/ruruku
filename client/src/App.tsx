import * as React from 'react';
import './App.css';

import logo from './logo.svg';
import { LoginForm } from './login-form';
import { Workspace } from './workspace';

interface AppState {
    name: string
    connecting: boolean
    socket?: WebSocket
}

class App extends React.Component<{}, AppState> {

    constructor(props: {}) {
        super(props);

        this.connect = this.connect.bind(this);
    }

    public componentWillMount() {
        this.setState({
            name: "none"
        });
    }

    public render() {
        let loginOrWorkspace;
        if (this.state.socket) {
            loginOrWorkspace = <Workspace name={this.state.name} ws={this.state.socket} />;
        } else if (this.state.connecting) {
            loginOrWorkspace = <div>Connecting ...</div>
        } else {
            loginOrWorkspace = <LoginForm handleSubmit={this.connect} />
        }

        return (
            <div className="app">
                <header>
                    <img src={logo} className="app-logo" alt="logo" /><h1 className="title">RURUKU</h1>
                </header>
                <div className="body">
                    { loginOrWorkspace }
                </div>
            </div>
        );
    }

    protected connect(name: string) {
        this.setState({ name, connecting: true });

        let protocol = 'ws';
        if (window.location.protocol === 'https:') {
            protocol = 'wss';
        }
        const ws = new WebSocket(`${protocol}://${window.location.host}/ws`);
        ws.onclose = (ev: CloseEvent) => this.setState({ socket: undefined, connecting: false });
        ws.onopen = () => this.setState({ socket: ws, connecting: false });
        ws.onerror = (err) => {
            console.log(err);
            alert(err);
            this.setState({ socket: ws, connecting: false });
        };
    }

}

export default App;
