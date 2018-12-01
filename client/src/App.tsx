import * as React from 'react';
import './App.css';
import 'semantic-ui-css/semantic.min.css';

import logo from './logo.svg';
// import { LoginForm } from './login-form';
// import { Workspace } from './workspace';
import { SessionList } from './session-list';
import { Participant } from './types/participant';
import { Workspace } from './workspace';

type AppMode = "select-session" | "in-session";

interface AppState {
    mode: AppMode
    participant?: Participant
    error?: string
}

class App extends React.Component<{}, AppState> {

    constructor(props: {}) {
        super(props);

        this.joinedSession = this.joinedSession.bind(this);
        this.state = { mode: "select-session" };
    }

    public render() {
        const body = this.state.participant && this.state.mode === "in-session"
            ? <Workspace participant={this.state.participant} />
            : <SessionList onJoin={this.joinedSession} />;

        return (
            <div className={"app " + this.state.mode}>
                <header id="header">
                    <img src={logo} className="app-logo" alt="logo" />
                </header>
                <div className="body">
                    {body}
                </div>
            </div>
        );
    }

    protected joinedSession(participant: Participant) {
        this.setState({ participant, mode: "in-session" });
    }

}

export default App;
