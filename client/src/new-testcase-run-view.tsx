import { TestCaseResult, TestCase } from '../../protocol/protocol';
import * as React from 'react';
import { Card, Button, ButtonOr, Form, Label, Dropdown, TextArea, DropdownProps, TextAreaProps } from 'semantic-ui-react';

export interface NewTestcaseRunViewProps {
    testcase: TestCase;
    onSubmit: (testcase: TestCase, result: TestCaseResult, comment: string) => void;
    onCancel: () => void;
}

interface NewTestcaseRunViewState {
    result: TestCaseResult;
    comment: string;
}

export class NewTestcaseRunView extends React.Component<NewTestcaseRunViewProps, NewTestcaseRunViewState> {

    constructor(props: NewTestcaseRunViewProps) {
        super(props);
        this.onCancel = this.onCancel.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
        this.onResultChange = this.onResultChange.bind(this);
        this.updateComment = this.updateComment.bind(this);
    }

    public componentWillMount() {
        this.setState({
            result: "undecided",
            comment: ""
        });
    }

    public render() {
        const tc = this.props.testcase;

        return <Card>
            <Card.Content>
                <Card.Header>{tc.name}</Card.Header>
                <Card.Meta>{tc.group} / {tc.id}</Card.Meta>
            </Card.Content>
            <Card.Content>
                <Form>
                    <Form.Field>
                        <Label>Result</Label>
                        <Dropdown
                            defaultValue="undecided"
                            options={[
                                { value: "passed", content: "Passed", icon: "check" },
                                { value: "undecided", content: "Undecided", icon: "question" },
                                { value: "failed", content: "Failed", icon: "times" }
                            ]}
                            search={true}
                            onChange={this.onResultChange}
                            value={this.state.result} />
                    </Form.Field>
                    <Form.Field>
                        <Label>Comment</Label>
                        <TextArea value={this.state.comment} onChange={this.updateComment} />
                    </Form.Field>
                </Form>
            </Card.Content>
            <Card.Content extra={true}>
                <ButtonOr>
                    <Button primary={true} color="green" onClick={this.onSubmit}>Submit</Button>
                    <Button basic={true} color="red" onClick={this.onCancel}>Cancel</Button>
                </ButtonOr>
            </Card.Content>
        </Card>;
    }

    protected onSubmit() {
        this.props.onSubmit(this.props.testcase, this.state.result, this.state.comment);
    }

    protected onCancel() {
        this.props.onCancel();
    }

    protected onResultChange(evt: React.SyntheticEvent, props: DropdownProps) {
        this.setState({ result: props.value as TestCaseResult });
    }

    protected updateComment(evt: React.SyntheticEvent, props: TextAreaProps) {
        this.setState({ comment: props.value as string });
    }

}