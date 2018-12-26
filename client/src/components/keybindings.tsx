import 'mousetrap';
import * as React from 'react';
import { Modal, Button, Icon, Table } from 'semantic-ui-react';

export class Keybindings {
    protected readonly bindings: Binding[][] = [[]];

    public push(...bindings: Binding[]) {
        this.bindings.push(bindings);
        this.bindings.forEach(bs => bs.forEach(b => Mousetrap.bind(b.keys, b.handler, 'keydown')));
    }

    public pop() {
        const bss = this.bindings.splice(this.bindings.length - 1, 1);
        bss.forEach(bs => bs.forEach(b => Mousetrap.unbind(b.keys, 'keydown')));
    }

    public getBindings(): Binding[] {
        return this.bindings.reduce((prev, curr) => ([] as Binding[]).concat(prev).concat(curr));
    }

}

export type BindingHandler = (e: ExtendedKeyboardEvent, combo: string) => void;

export interface Binding {
    keys: string | string[];
    description: string;
    handler: BindingHandler;
}

export interface KeybindingHelpDialogProps {
    bindings: Keybindings;
    open: boolean;
    closeDialog: () => void;
}

export class KeybindingHelpDialog extends React.Component<KeybindingHelpDialogProps, {}> {

    public render() {
        const rows = this.props.bindings.getBindings().map(b => (
            <Table.Row>
                <Table.Cell>{b.keys}</Table.Cell>
                <Table.Cell>{b.description}</Table.Cell>
            </Table.Row>
        ))
        return (
            <Modal open={this.props.open}>
                <Modal.Header>
                    Help
                </Modal.Header>
                <Modal.Content>
                    <h2>Available keybindings</h2>
                    <Table>
                        <Table.Header>
                            <Table.Row>
                                <Table.HeaderCell>Key</Table.HeaderCell>
                                <Table.HeaderCell>Description</Table.HeaderCell>
                            </Table.Row>
                        </Table.Header>
                        <Table.Body>
                            {rows}
                        </Table.Body>
                    </Table>
                </Modal.Content>
                <Modal.Actions>
                    <Button color='green' inverted={true} onClick={this.props.closeDialog}><Icon name='checkmark' /> Ok</Button>
                </Modal.Actions>
            </Modal>
        )
    }

}
