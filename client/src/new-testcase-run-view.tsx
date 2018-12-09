import * as React from 'react';
import { Card, Button, ButtonOr, Form, Label, Dropdown, TextArea, DropdownProps, TextAreaProps, ButtonGroup } from 'semantic-ui-react';
import { TestcaseDetailViewBody } from './testcase-detail-view-body';
import { Testcase, TestRunState, TestcaseStatus, TestcaseRunResult } from './api/v1/session_pb';
import { Participant } from './types/participant';

export interface NewTestcaseRunViewProps {
    testcase: TestcaseStatus;
    participant: Participant;
    result?: TestRunState;
    onSubmit: (testCase: Testcase, participant: Participant, result: TestRunState, comment: string) => void;
    onClose: () => void;
}

interface NewTestcaseRunViewState {
    result: TestRunState;
    comment: string;
}

function findContributionByParticipant(tc: TestcaseStatus, participant: string): TestcaseRunResult | undefined {
    return tc.getResultList().find(r => {
        const p = r.getParticipant();
        if (p) {
            if (p.getName() === participant) {
                return true;
            }
        }
        return false;
    })
}

export class NewTestcaseRunView extends React.Component<NewTestcaseRunViewProps, NewTestcaseRunViewState> {
    protected focusElement: any | undefined;

    constructor(props: NewTestcaseRunViewProps) {
        super(props);

        const ourContribution = findContributionByParticipant(props.testcase, props.participant.name);
        if (ourContribution !== undefined) {
            this.state = {
                result: this.props.result !== undefined ? this.props.result : ourContribution.getState(),
                comment: ourContribution.getComment()
            };
        } else {
            this.state = {
                result: this.props.result !== undefined ? this.props.result : TestRunState.UNDECIDED,
                comment: ""
            };
        }

        this.setFocusElement = this.setFocusElement.bind(this);
        this.onCancel = this.onCancel.bind(this);
        this.onSubmit = this.onSubmit.bind(this);
        this.onResultChange = this.onResultChange.bind(this);
        this.updateComment = this.updateComment.bind(this);
    }

    public componentDidMount(){
        if (this.focusElement) {
            // this.focusElement.focus();
        }
    }

    public render() {
        const tcs = this.props.testcase;
        const tc = tcs.getCase()!;
        const resultOptions = [
            { value: TestRunState.PASSED, text: "Passed", icon: "check" },
            { value: TestRunState.UNDECIDED, text: "Undecided", icon: "question" },
            { value: TestRunState.FAILED, text: "Failed", icon: "times" }
        ];

        return <Card>
            <Card.Content>
                <Card.Header>{tc.getName()}</Card.Header>
                <Card.Meta>{tc.getGroup()} / {tc.getId()}</Card.Meta>
            </Card.Content>
            <Card.Content>
                <Form>
                    <Form.Field>
                        <Label>Result</Label>
                        <Dropdown
                            options={resultOptions}
                            onChange={this.onResultChange}
                            value={this.state.result} />
                    </Form.Field>
                    <Form.Field>
                        <Label>Comment</Label>
                        <TextArea value={this.state.comment} onChange={this.updateComment} ref={this.setFocusElement} />
                    </Form.Field>
                </Form>
            </Card.Content>
            <Card.Content extra={true}>
                <ButtonGroup>
                    <Button positive={true} onClick={this.onSubmit} fluid={true}>Submit</Button>
                    <ButtonOr />
                    <Button onClick={this.onCancel}>Cancel</Button>
                </ButtonGroup>
            </Card.Content>
            <TestcaseDetailViewBody tcs={tcs} />
        </Card>;
    }

    protected setFocusElement(element: any) {
        this.focusElement = element;
    }

    protected onSubmit() {
        this.props.onSubmit(this.props.testcase.getCase()!, this.props.participant, this.state.result, this.state.comment);
        this.props.onClose();
    }

    protected onCancel() {
        this.props.onClose();
    }

    protected onResultChange(evt: React.SyntheticEvent, props: DropdownProps) {
        this.setState({ result: props.value as TestRunState });
    }

    protected updateComment(evt: React.SyntheticEvent, props: TextAreaProps) {
        this.setState({ comment: props.value as string });
    }

}