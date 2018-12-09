import * as React from 'react';
import { Card, Button } from 'semantic-ui-react';
import { Converter } from 'showdown';
import { TestcaseDetailViewBody } from './testcase-detail-view-body';
import { TestcaseStatus } from './api/v1/session_pb';

export interface TestcaseDetailViewProps {
    testcase: TestcaseStatus
    onClose: () => void
}

export class TestcaseDetailView extends React.Component<TestcaseDetailViewProps, {}> {

    constructor(pros: TestcaseDetailViewProps) {
        super(pros);
        this.onClose = this.onClose.bind(this);
    }

    public render() {
        const tcs = this.props.testcase;
        const tc = tcs.getCase()!;

        const markdown = new Converter({
            headerLevelStart: 4
        });

        return <Card>
            <Card.Content>
                <Card.Header>{tc.getName()}</Card.Header>
                <Card.Meta>{tc.getGroup()} / {tc.getId()}</Card.Meta>
                <Card.Description><div dangerouslySetInnerHTML={{__html: markdown.makeHtml(tc.getDescription())}} /></Card.Description>
            </Card.Content>
            <TestcaseDetailViewBody tcs={tcs} />
            <Card.Content extra={true}>
                <Button basic={true} color="red" onClick={this.onClose}>Close</Button>
            </Card.Content>
        </Card>;
    }

    protected onClose() {
        this.props.onClose();
    }

}
