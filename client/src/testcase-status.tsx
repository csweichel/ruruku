import { TestCase, TestCaseRun, TestParticipant } from '../../protocol/protocol';
import * as React from 'react';
import './testcase-status.css';
import { Icon } from 'semantic-ui-react';

export interface TestCaseParticipant {
    run?: TestCaseRun
    participant?: TestParticipant
}

export interface TestCaseStatusViewProps {
    case: TestCase
    runs: TestCaseParticipant[]
}

export class TestCaseStatusView extends React.Component<TestCaseStatusViewProps, {}> {

    public render() {
        const remainingCellCount = Math.max(this.props.case.minTesterCount, this.props.runs.length) - this.props.runs.length;
        const content = [];
        for(const r of this.props.runs) {
            if (r.run) {
                if (r.run.result === "passed") {
                    content.push(<Icon name="check" className="passed" key={content.length}  about="Passed" />);
                } else if (r.run.result === "failed") {
                    content.push(<Icon name="times" className="failed" key={content.length} about="Failed" />);
                } else {
                    content.push(<Icon name="circle outline" className="undecided" key={content.length} about="Undecided" />);
                }
            } else {
                content.push(<Icon name="user" key={content.length} about={r.participant!.name} />);
            }
        }
        for(let i = 0; i < remainingCellCount; i++) {
            content.push(<Icon name="user outline" key={content.length} about="Unclaimed" />);
        }
        return <div className="case-status-view">{content}</div>
    }

}