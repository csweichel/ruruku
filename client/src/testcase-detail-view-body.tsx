import * as React from 'react';
import { TestCaseRun, TestCase } from '../../protocol/protocol';
import { Converter } from 'showdown';
import { Card, Icon, Feed } from 'semantic-ui-react';

export class TestcaseDetailViewBody extends React.Component<{ tc: TestCase, runs?: TestCaseRun[] }, {}> {

    public render() {
        const tc = this.props.tc;
        const markdown = new Converter({
            headerLevelStart: 4
        });

        const mustPass = tc.mustPass ? (
            <Card.Content extra={true} key="mustpass">
                must pass
            </Card.Content>
        ) : undefined;

        const steps = tc.steps ? (
            <Card.Content key="steps">
                <h3>Steps</h3>
                <div dangerouslySetInnerHTML={{__html: markdown.makeHtml(tc.steps)}} />
            </Card.Content>
        ) : undefined;

        let runs: JSX.Element | undefined;
        if (this.props.runs && this.props.runs.length > 0) {
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
            runs = (
                <Card.Content key="runs">
                    <Feed>{allRuns}</Feed>
                </Card.Content>
            );
        }

        return [ steps, runs, mustPass ].filter(e => !!e);
    }

}