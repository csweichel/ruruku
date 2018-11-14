import * as React from 'react';

export interface LoginFormProps {
    handleSubmit (name: string): void
}

interface LoginFormState {
    name: string
}

export class LoginForm extends React.Component<LoginFormProps, LoginFormState> {

    constructor(props: LoginFormProps) {
        super(props);

        this.onChange = this.onChange.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
    }

    public componentWillMount() {
        this.setState({
            name: ""
        });
    }

    public render() {
        return (
            <form onSubmit={this.onSubmit}>
                <input type="text" placeholder="Choose a name..." value={this.state.name} onChange={this.onChange} />
                <input type="submit" value="Join" disabled={!this.state.name} />
            </form>
        );
    }

    protected onSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        if(this.state.name) {
            this.props.handleSubmit(this.state.name);
        }
    }

    protected onChange(e: React.FormEvent<HTMLInputElement>) {
        this.setState({ name: (e.target as HTMLInputElement).value });
    }

}