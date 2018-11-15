import { TestCase, TestCaseRun } from '../../protocol/protocol';
import * as React from 'react';
import './testcase-status.css';

export interface TestCaseStatusViewProps {
    case: TestCase
    runs: TestCaseRun[]
}

export class TestCaseStatusView extends React.Component<TestCaseStatusViewProps, {}> {

    public render() {
        const remainingCellCount = Math.max(this.props.case.minTesterCount, this.props.runs.length) - this.props.runs.length;
        const content = [];
        for(const r of this.props.runs) {
            content.push(<div className={r.result} key={content.length} />);
        }
        for(let i = 0; i < remainingCellCount; i++) {
            content.push(<div className="unassigned" key={content.length} />);
        }
        return <div className="case-status-view">{content}</div>
    }

}