import * as React from 'react';
import { TestSuite, TestRun, TestParticipant } from '../../protocol/protocol';
import { Table, Checkbox } from 'semantic-ui-react';
import { TestCaseStatusView } from './testcase-status';

export interface TestplanViewProps {
    suite: TestSuite
    run: TestRun
    participant: TestParticipant
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
        const matchedCases = this.props.suite.cases.map(cse => {
            const runs = (this.props.run.cases || []).filter(cser => cser.caseName === cse.name && cser.caseGroup === cse.group);
            return { case: cse, runs };
        });

        return matchedCases.map((mc, idx) => {
            return <Table.Row key={idx}>
                <Table.Cell><Checkbox checked={this.isClaimed(mc.case.name)} /></Table.Cell>
                <Table.Cell>{mc.case.group}</Table.Cell>
                <Table.Cell>{mc.case.name}</Table.Cell>
                <Table.Cell><TestCaseStatusView case={mc.case} runs={mc.runs} /></Table.Cell>
            </Table.Row>
        });
    }

    protected isClaimed(testCase: string): boolean {
        return this.props.participant.claimedCases.indexOf(testCase) > -1;
    }

}