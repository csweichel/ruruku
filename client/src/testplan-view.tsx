import * as React from 'react';
import { TestSuite, TestRun } from '../../protocol/protocol';
import { Table } from 'semantic-ui-react'

export interface TestplanViewProps {
    suite: TestSuite
    run: TestRun
}

export class TestplanView extends React.Component<TestplanViewProps, {}> {

    public render() {
        return <Table celled={true} compact={true} definition={true}>
            <Table.Header>
                <Table.Row>
                    <Table.Cell>Group</Table.Cell>
                    <Table.Cell>Name</Table.Cell>
                    <Table.Cell>Tester</Table.Cell>
                </Table.Row>
            </Table.Header>
        </Table>
    }

}