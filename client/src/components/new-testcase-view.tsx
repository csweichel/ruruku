import * as React from 'react';
import { Button, ButtonOr, Form, ButtonGroup, Input, FormTextArea, Checkbox, InputProps, CheckboxProps } from 'semantic-ui-react';
import { Testcase } from '../api/v1/session_pb';
import { AppStateContent, getClient } from 'src/types/app-state';

export interface NewTestcaseViewProps {
    appState: AppStateContent;
    testcase?: Testcase;
    onClose: () => void;
}

interface NewTestcaseViewState {
    id?: string;
    group?: string;
    name?: string;
    mustPass?: boolean;
    minTesterCount?: number;
    description?: string;
    steps?: string;
}

export class NewTestcaseView extends React.Component<NewTestcaseViewProps, NewTestcaseViewState> {
    protected focusElement: any | undefined;

    constructor(props: NewTestcaseViewProps) {
        super(props);

        if (this.props.testcase) {
            this.state = {
                id: this.props.testcase.getId(),
                group: this.props.testcase.getGroup(),
                name: this.props.testcase.getName(),
                mustPass: this.props.testcase.getMustpass(),
                minTesterCount: this.props.testcase.getMintestercount(),
                description: this.props.testcase.getDescription(),
                steps: this.props.testcase.getSteps(),
            }
        } else {
            this.state = {
                minTesterCount: 1,
                mustPass: true
            }
        }

        this.setFocusElement = this.setFocusElement.bind(this);
        this.onCancel = this.onCancel.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
        this.onChange = this.onChange.bind(this);
    }

    public componentDidMount(){
        if (this.focusElement) {
            // this.focusElement.focus();
        }
    }

    public render() {
        return (
            <Form>
                <Form.Field>
                    <label>ID</label>
                    <Input disabled={!!this.props.testcase} type="text" value={this.state.id} placeholder="name" onChange={this.onChange.bind(this, "id")} />
                </Form.Field>
                <Form.Field>
                    <label>Group</label>
                    <Input type="text" value={this.state.group} placeholder="name" onChange={this.onChange.bind(this, "group")} />
                </Form.Field>
                <Form.Field>
                    <label>Name</label>
                    <Input type="text" value={this.state.name} placeholder="name" onChange={this.onChange.bind(this, "name")} />
                </Form.Field>
                <Form.Field>
                    <label>Minimum Tester Count</label>
                    <Input type="text" value={this.state.minTesterCount} placeholder="name" onChange={this.onChange.bind(this, "minTesterCount")} />
                </Form.Field>
                <Form.Field>
                    <Checkbox toggle={true} label="Must pass" checked={this.state.mustPass} onChange={this.onChange.bind(this, "mustPass")} />
                </Form.Field>
                <Form.Field>
                    <label>Description</label>
                    <FormTextArea content={this.state.description} onChange={this.onChange.bind(this, "description")} />
                </Form.Field>
                <Form.Field>
                    <label>Steps</label>
                    <FormTextArea content={this.state.steps} onChange={this.onChange.bind(this, "steps")} />
                </Form.Field>

                <ButtonGroup>
                    <Button positive={true} onClick={this.onSubmit} fluid={true}>Submit</Button>
                    <ButtonOr />
                    <Button onClick={this.onCancel}>Cancel</Button>
                </ButtonGroup>
            </Form>
        );
    }

    protected setFocusElement(element: any) {
        this.focusElement = element;
    }

    protected async onSubmit() {
        const tc = new Testcase();
        tc.setId(this.state.id || "");
        tc.setGroup(this.state.group || "");
        tc.setName(this.state.name || "");
        tc.setMintestercount(this.state.minTesterCount || 0);
        tc.setMustpass(!!this.state.mustPass);
        tc.setDescription(this.state.description || "");
        tc.setSteps(this.state.steps || "");

        if (!!this.props.testcase) {
            tc.setId(this.props.testcase.getId());
            try {
                await getClient(this.props.appState).modifyTestcase(tc);
                this.props.onClose();
            } catch (err) {
                this.props.appState.setError(err);
            }
        } else {
            try {
                await getClient(this.props.appState).addTestcase(tc);
                this.props.onClose();
            } catch (err) {
                this.props.appState.setError(err);
            }
        }
    }

    protected onCancel() {
        this.props.onClose();
    }

    protected onChange(field: string, evt: React.SyntheticEvent, props: any) {
        if (field === "minTesterCount") {
            const n = parseInt((props as InputProps).value, 10);
            if (!Number.isNaN(n)) {
                this.setState({ minTesterCount: n });
            }
        } else if (field === "mustPass") {
            this.setState({ mustPass: !!(props as CheckboxProps).value });
        } else {
            const val = (props as InputProps).value;
            const s: any = {};
            s[field] = val;
            this.setState(s);
        }
    }

}