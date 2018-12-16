import * as React from 'react';
import './App.css';
import 'semantic-ui-css/semantic.min.css';

import logo from './logo.svg';
import { AppStateContent } from './types/app-state';
import { LoginForm } from './components/login-form';
import { SessionSelector } from './components/session-selector';
import { Workspace } from './components/workspace';

class App extends React.Component<{}, AppStateContent> {

    constructor(props: {}) {
        super(props);

        let errorTimeout: NodeJS.Timeout | undefined;
        this.state = {
            setError: error => {
                this.setState({ error });
                if (error) {
                    if (errorTimeout) {
                        clearTimeout(errorTimeout);
                    }
                    errorTimeout = setTimeout(() => this.setState({ error: undefined }), 5000);
                }
            }
        };

        this.onLogin = this.onLogin.bind(this);
    }

    public render() {
        let className = "login";
        let content: JSX.Element | undefined;
        if (!this.state.token) {
            content = <LoginForm onLogin={this.onLogin} appState={this.state} />;
        } else if (!this.state.session) {
            const selectSession = (session: string) => this.setState({ session });
            content = <SessionSelector onSelect={selectSession} appState={this.state} />;
        } else {
            className = "workspace";
            content = <Workspace appState={this.state} />
        }

        return (
            <div className={`app ${className}`}>
                <header id="header">
                    <img src={logo} className="app-logo" alt="logo" />
                </header>
                <div className="body">
                    {content}
                </div>
            </div>
        );
    }

    protected onLogin(token: string) {
        this.setState({ token });
    }

}

export default App;
