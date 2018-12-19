import { AppStateContent, getClient } from 'src/types/app-state';
import * as React from 'react';
import { Form, Label, Input, Button, Segment, Message, Modal, Icon } from 'semantic-ui-react';

export interface ChangePasswordFormProps {
    appState: AppStateContent;
    onClose: () => void;
}

export interface ChangePasswordFormState {
    oldPassword?: string;
    newPassword?: string;
    repeatPassword?: string;
    done?: boolean;
}

export class ChangePasswordForm extends React.Component<ChangePasswordFormProps, ChangePasswordFormState> {

    constructor(props: ChangePasswordFormProps) {
        super(props);
        this.state = {};

        this.onSubmit = this.onSubmit.bind(this);
    }

    public render() {
        let error: JSX.Element | undefined;
        if (this.props.appState.error) {
            error = <Message error={true}>{this.props.appState.error}</Message>
        }
        return (
            <Segment>
                {error}
                <Form onSubmit={this.onSubmit}>
                    <Form.Field>
                        <Label>Old Password</Label>
                        <Input type="password" placeholder="Old Password" value={this.state.oldPassword} onChange={this.onChange.bind(this, "oldPassword")} />
                    </Form.Field>
                    <Form.Field>
                        <Label>New Password</Label>
                        <Input type="password" placeholder="New Password" value={this.state.newPassword} onChange={this.onChange.bind(this, "newPassword")} />
                    </Form.Field>
                    <Form.Field>
                        <Input type="password" placeholder="Repeat Password" value={this.state.repeatPassword} onChange={this.onChange.bind(this, "repeatPassword")} />
                    </Form.Field>
                    <Form.Field>
                        <Button basic={true} color='red' onClick={this.props.onClose}>Cancel</Button>
                        <Button color='green' type="submit" primary={true} disabled={!this.state.newPassword || this.state.newPassword.length <= 4 || this.state.newPassword !== this.state.repeatPassword}>Change</Button>
                    </Form.Field>
                </Form>
                <Modal open={!!this.state.done}>
                    <Modal.Content>
                        <p>Your password was changed.</p>
                    </Modal.Content>
                    <Modal.Actions>
                        <Button color='green' onClick={this.props.onClose}><Icon name='checkmark' /> Ok</Button>
                    </Modal.Actions>
                </Modal>
            </Segment>
        )
    }

    protected async onSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();

        try {
            await getClient(this.props.appState).changePassword(this.state.oldPassword!, this.state.newPassword!);
            this.setState({ done: true });
        } catch(err) {
            this.props.appState.setError(err);
        }
    }

    protected onChange(field: string, le: React.SyntheticEvent, e: HTMLInputElement) {
        if (field === 'oldPassword') {
            this.setState({ oldPassword: e.value });
        }
        if (field === 'newPassword') {
            this.setState({ newPassword: e.value });
        }
        if (field === 'repeatPassword') {
            this.setState({ repeatPassword: e.value });
        }
    }

}