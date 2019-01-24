import * as React from 'react';
import { Input, Form, Segment, Button, Message } from 'semantic-ui-react';

import "./login-form.css";
import { AppStateContent, getClient } from 'src/types/app-state';

export interface LoginFormProps {
    appState: AppStateContent;
    onLogin: (token: string) => void
}

interface LoginFormState {
    name: string
    password: string
}

export class LoginForm extends React.Component<LoginFormProps, LoginFormState> {

    constructor(props: LoginFormProps) {
        super(props);

        this.onChange = this.onChange.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
    }

    public componentWillMount() {
        this.setState({
            name: "",
            password: ""
        });
    }

    public render() {
        let error: JSX.Element | undefined;
        if (this.props.appState.error) {
            error = <Message error={true}>{this.props.appState.error}</Message>;
        }
        return (
            <Segment id="login-form">
                {error}
                <Form onSubmit={this.onSubmit}>
                    <Form.Field>
                        <Input type="text" placeholder="Username" value={this.state.name} onChange={this.onChange.bind("username")} />
                    </Form.Field>
                    <Form.Field>
                        <Input type="password" placeholder="Password" value={this.state.password} onChange={this.onChange.bind("password")} />
                    </Form.Field>
                    <Form.Field>
                        <Button type="submit" primary={true} disabled={!this.state.name || !this.state.password}>Login</Button>
                    </Form.Field>
                </Form>
            </Segment>
        );
    }

    protected async onSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();

        try {
            const token = await getClient(this.props.appState).login(this.state.name, this.state.password);
            this.props.onLogin(token);
        } catch(err) {
            this.props.appState.setError(err);
        }
    }

    protected onChange(field: string, e: HTMLInputElement) {
        const value = e.value;
        if (e.type === 'password') {
            this.setState({ password: value });
        } else {
            this.setState({ name: value });
        }
    }

}