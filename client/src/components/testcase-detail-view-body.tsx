import * as React from 'react';
import { Converter } from 'showdown';
import { Card, Icon, Feed, Table } from 'semantic-ui-react';
import { TestcaseStatus, TestRunState } from '../api/v1/session_pb';

export class TestcaseDetailViewBody extends React.Component<{ tcs: TestcaseStatus }, {}> {

    public render() {
        const tcs = this.props.tcs;
        const tc = tcs.getCase()!;
        const markdown = new Converter({
            headerLevelStart: 4
        });

        const mustPass = tc.getMustpass() ? (
            <Card.Content extra={true} key="mustpass">
                must pass
            </Card.Content>
        ) : undefined;

        const steps = tc.getSteps() ? (
            <Card.Content key="steps">
                <h3>Steps</h3>
                <div dangerouslySetInnerHTML={{__html: markdown.makeHtml(tc.getSteps())}} />
            </Card.Content>
        ) : undefined;

        const annotationMap = tc.getAnnotationsMap();
        const annotations = !annotationMap.keys().next().done ? (
            <Card.Content key="annotations">
                <h3>Annotations</h3>
                <Table definition={true}>
                    { annotationMap.getEntryList().map(kv => <tr key={kv[0]}><td>{kv[0]}</td><td>{kv[1]}</td></tr>) }
                </Table>
            </Card.Content>
        ) : undefined;

        let runs: JSX.Element | undefined;
        const results = tcs.getResultList();
        if (results && results.length > 0) {
            const allRuns = results.map((r, idx) => {
                let icon = <Icon name="circle outline" />;
                if (r.getState() === TestRunState.PASSED) {
                    icon = <Icon name="check" />;
                } else if (r.getState() === TestRunState.FAILED) {
                    icon = <Icon name="times" />;
                } else {
                    icon = <Icon name="question" />;
                }

                return <Feed.Event key={idx}>
                    <Feed.Label>{icon}</Feed.Label>
                    <Feed.Content>
                        <Feed.Date>{r.getParticipant()!.getName()}</Feed.Date>
                        <div dangerouslySetInnerHTML={{__html: markdown.makeHtml(r.getComment())}} />
                    </Feed.Content>
                </Feed.Event>
            });
            runs = (
                <Card.Content key="runs">
                    <Feed>{allRuns}</Feed>
                </Card.Content>
            );
        }

        return [ steps, runs, mustPass, annotations ].filter(e => !!e);
    }

}