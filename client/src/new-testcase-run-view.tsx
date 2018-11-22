import { TestCaseResult, TestCase } from '../../protocol/protocol';
import * as React from 'react';
import { Card, Button, ButtonOr, Form, Label, Dropdown, TextArea, DropdownProps, TextAreaProps, ButtonGroup } from 'semantic-ui-react';

export interface NewTestcaseRunViewProps {
    testcase: TestCase;
    onSubmit: (testcase: TestCase, result: TestCaseResult, comment: string) => void;
    onClose: () => void;
    claimTestCase: (testcase: TestCase, claim: boolean) => void;
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
        this.onUnclaim = this.onUnclaim.bind(this);
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
                            options={[
                                { value: "passed", text: "Passed", icon: "check" },
                                { value: "undecided", text: "Undecided", icon: "question" },
                                { value: "failed", text: "Failed", icon: "times" }
                            ]}
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
                <ButtonGroup>
                    <Button positive={true} onClick={this.onSubmit}>Submit</Button>
                    <ButtonOr />
                    <Button onClick={this.onCancel}>Cancel</Button>
                </ButtonGroup>
            </Card.Content>
            <Card.Content extra={true}>
                <Button basic={true} onClick={this.onUnclaim} fluid={true}>Unclaim</Button>
            </Card.Content>
        </Card>;
    }

    protected onSubmit() {
        this.props.onSubmit(this.props.testcase, this.state.result, this.state.comment);
        this.props.onClose();
    }

    protected onCancel() {
        this.props.onClose();
    }

    protected onUnclaim() {
        this.props.claimTestCase(this.props.testcase, false);
        this.props.onClose();
    }

    protected onResultChange(evt: React.SyntheticEvent, props: DropdownProps) {
        this.setState({ result: props.value as TestCaseResult });
    }

    protected updateComment(evt: React.SyntheticEvent, props: TextAreaProps) {
        this.setState({ comment: props.value as string });
    }

}