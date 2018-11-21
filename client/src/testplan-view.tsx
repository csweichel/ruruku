import * as React from 'react';
import { TestSuite, TestRun, TestParticipant, TestCase, TestCaseRun } from '../../protocol/protocol';
import { Table, Button } from 'semantic-ui-react';
import { TestCaseStatusView, TestCaseParticipant } from './testcase-status';
import { TestcaseDetailView } from './testcase-detail-view';

export interface TestplanViewProps {
    suite: TestSuite
    run: TestRun
    participant: TestParticipant

    claimTestCase(testCase: TestCase, claim: boolean): void
    showDetails(content?: any): void
}

export class TestplanView extends React.Component<TestplanViewProps, {}> {

    public render() {
        const testcases = this.buildRows();
        return <Table celled={true} compact={true} sortable={true} fixed={true}>
            <Table.Header>
                <Table.Row>
                    <Table.HeaderCell />
                    <Table.HeaderCell>Group</Table.HeaderCell>
                    <Table.HeaderCell>Name</Table.HeaderCell>
                    <Table.HeaderCell>Testers</Table.HeaderCell>
                </Table.Row>
            </Table.Header>
            <Table.Body>
                {testcases}
            </Table.Body>
        </Table>
    }

    protected buildRows() {
        if (!this.props.suite) {
            return [];
        }

        const matchedCases = this.props.suite.cases.map(cse => {
            const runs = (this.props.run.cases || []).filter(cser => cser.case === cse.id && cser.caseGroup === cse.group);
            return { case: cse, runs };
        });

        return matchedCases.map(mc => {
            return <Table.Row key={mc.case.id}>
                <Table.Cell>
                    {this.getActions(mc.case)}
                </Table.Cell>
                <Table.Cell>{mc.case.group}</Table.Cell>
                <Table.Cell><a href="#" onClick={this.showDetails.bind(this, mc.case, mc.runs)}>{mc.case.name}</a></Table.Cell>
                <Table.Cell><TestCaseStatusView case={mc.case} runs={this.getRunsAndClaims(mc.case, mc.runs)} /></Table.Cell>
            </Table.Row>
        });
    }

    protected getActions(tc: TestCase) {
        if (this.isClaimed(tc)) {
            return [
                <Button label="Complete" icon="write square" key="complete" />,
                <Button label="Unclaim" basic={true} icon="minus circle" onClick={this.claim.bind(this, tc, false)} key="unclaim" />
            ]
        } else {
            return [
                <Button label="Claim" icon="plus circle" onClick={this.claim.bind(this, tc, true)} key="claim" />
            ]
        }
    }

    protected showDetails(cse: TestCase, runs: TestCaseRun[], evt: React.SyntheticEvent) {
        evt.preventDefault();
        this.props.showDetails(<TestcaseDetailView testcase={cse} runs={runs} onClose={this.props.showDetails} />);
    }

    protected getRunsAndClaims(cse: TestCase, runs: TestCaseRun[]): TestCaseParticipant[] {
        const participants = this.props.run.participants.filter(p => Object.keys(p.claimedCases).indexOf(`${cse.group}/${cse.id}`) > -1);

        return participants.map(p => {
            return {
                participant: p
            } as TestCaseParticipant
        }).concat(runs.map(r => {
            return {
                run: r
            } as TestCaseParticipant
        }));
    }

    protected isClaimed(cse: TestCase): boolean {
        return this.props.participant.claimedCases[`${cse.group}/${cse.id}`];
    }

    protected claim(cse: TestCase, claim: boolean, evt: React.SyntheticEvent): void {
        evt.preventDefault();
        this.props.claimTestCase(cse, claim);
    }

}