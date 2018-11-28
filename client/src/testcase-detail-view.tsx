import * as React from 'react';
import { TestCase, TestCaseRun } from '../../protocol/protocol';
import { Card, Button } from 'semantic-ui-react';
import { Converter } from 'showdown';
import { TestcaseDetailViewBody } from './testcase-detail-view-body';

export interface TestcaseDetailViewProps {
    testcase: TestCase
    runs: TestCaseRun[]
    onClose: () => void
}

export class TestcaseDetailView extends React.Component<TestcaseDetailViewProps, {}> {

    constructor(pros: TestcaseDetailViewProps) {
        super(pros);
        this.onClose = this.onClose.bind(this);
    }

    public render() {
        const tc = this.props.testcase;

        const markdown = new Converter({
            headerLevelStart: 4
        });

        return <Card>
            <Card.Content>
                <Card.Header>{tc.name}</Card.Header>
                <Card.Meta>{tc.group} / {tc.id}</Card.Meta>
                <Card.Description><div dangerouslySetInnerHTML={{__html: markdown.makeHtml(tc.description)}} /></Card.Description>
            </Card.Content>
            <TestcaseDetailViewBody tc={tc} runs={this.props.runs} />
            <Card.Content extra={true}>
                <Button basic={true} color="red" onClick={this.onClose}>Close</Button>
            </Card.Content>
        </Card>;
    }

    protected onClose() {
        this.props.onClose();
    }

}
