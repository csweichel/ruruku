import * as React from 'react';
import './App.css';
import 'semantic-ui-css/semantic.min.css';

import logo from './logo.svg';
import { AppStateContent } from './types/app-state';
import { LoginForm } from './components/login-form';
import { SessionSelector } from './components/session-selector';
import { Workspace } from './components/workspace';
import { MiniEventEmitter } from './types/mini-event-emitter';
import { Modal, Button, Icon } from 'semantic-ui-react';
import { Keybindings, KeybindingHelpDialog } from './components/keybindings';

export interface AppProps {
    reloadRequest: MiniEventEmitter<boolean>;
}

interface AppPrivateState {
    shouldReload?: boolean;
    showKeybindingHelp?: boolean;
}

class App extends React.Component<AppProps, AppStateContent & AppPrivateState> {

    constructor(props: AppProps) {
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
            },
            logout: () => this.setState({ token: undefined }),
            resetSession: () => this.setState({ session: undefined }),
            keybindings: new Keybindings(),
        };

        this.onLogin = this.onLogin.bind(this);
        this.ignoreReloadRequest = this.ignoreReloadRequest.bind(this);
    }

    public componentWillMount() {
        this.props.reloadRequest.subscribe(() => {
            this.setState({ shouldReload: true });
        });
        this.state.keybindings.push({
            keys: 'h', description: 'Show help dialog', handler: () => this.showKeyboardHelpDialog(true)
        });
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
                <Modal basic={true} open={!!this.state.shouldReload}>
                    <Modal.Content>
                        <p>A new version of Ruruku is available. Please reload this page.</p>
                    </Modal.Content>
                    <Modal.Actions>
                        <Button basic={true} color='red' inverted={true} onClick={this.ignoreReloadRequest}>Ignore</Button>
                        <Button color='green' inverted={true} onClick={this.reloadPage}><Icon name='checkmark' /> Do it</Button>
                    </Modal.Actions>
                </Modal>
                <KeybindingHelpDialog bindings={this.state.keybindings} open={!!this.state.showKeybindingHelp} closeDialog={this.showKeyboardHelpDialog.bind(this, false)} />
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

    protected ignoreReloadRequest() {
        this.setState({ shouldReload: false });
    }

    protected reloadPage() {
        window.location.reload();
    }

    protected showKeyboardHelpDialog(open: boolean) {
        this.setState({ showKeybindingHelp: open });
    }

}

export default App;
