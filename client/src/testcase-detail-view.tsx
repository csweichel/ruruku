import * as React from 'react';
import { TestCase, TestCaseRun } from '../../protocol/protocol';
import { Card, Feed, Icon, Button } from 'semantic-ui-react';
import { Converter } from 'showdown';

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

        const mustPass = tc.mustPass ? (
            <Card.Content extra={true}>
                must pass
            </Card.Content>
        ) : undefined;

        const steps = tc.steps ? (
            <Card.Content>
                <h3>Steps</h3>
                <div dangerouslySetInnerHTML={{__html: markdown.makeHtml(tc.steps)}} />
            </Card.Content>
        ) : undefined;

        const allRuns = this.props.runs.map((r, idx) => {
            let icon = <Icon name="circle outline" />;
            if (r.result === 'passed') {
                icon = <Icon name="check" />;
            } else if (r.result === 'failed') {
                icon = <Icon name="times" />;
            } else {
                icon = <Icon name="question" />;
            }

            return <Feed.Event key={idx}>
                <Feed.Label>{icon}</Feed.Label>
                <Feed.Content>
                    <Feed.Date>{r.tester}</Feed.Date>
                    <div dangerouslySetInnerHTML={{__html: markdown.makeHtml(r.comment)}} />
                </Feed.Content>
            </Feed.Event>
        });
        const runs = allRuns.length > 0 ? (
            <Card.Content>
                <Feed>{allRuns}</Feed>
            </Card.Content>
        ) : undefined;
        return <Card>
            <Card.Content>
                <Card.Header>{tc.name}</Card.Header>
                <Card.Meta>{tc.group} / {tc.id}</Card.Meta>
                <Card.Description><div dangerouslySetInnerHTML={{__html: markdown.makeHtml(tc.description)}} /></Card.Description>
            </Card.Content>
            {steps}
            {runs}
            {mustPass}
            <Card.Content extra={true}>
                <Button basic={true} color="red" onClick={this.onClose}>Close</Button>
            </Card.Content>
        </Card>;
    }

    protected onClose() {
        this.props.onClose();
    }

}