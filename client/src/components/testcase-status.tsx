import * as React from 'react';
import './testcase-status.css';
import { Icon } from 'semantic-ui-react';
import { TestcaseStatus, TestRunState } from '../api/v1/session_pb';

export interface TestCaseStatusViewProps {
    case: TestcaseStatus
}

export class TestCaseStatusView extends React.Component<TestCaseStatusViewProps, {}> {

    public render() {
        const cse = this.props.case.getCase()!;
        const results = this.props.case.getResultList() || [];
        const claims = this.props.case.getClaimList() || [];
        const claimAndContributionCount = results.length + claims.length;
        const remainingCellCount = Math.max(cse.getMintestercount(), claimAndContributionCount) - claimAndContributionCount;
        const resolvedClaims = new Map<string, boolean>();

        const content = [];
        for(const r of results) {
            const participantName = r.getParticipant()!.getName();
            resolvedClaims.set(participantName, true);

            const state = r.getState();
            if (state === TestRunState.PASSED) {
                content.push(<Icon name="check" className="passed" key={content.length}  about="Passed" />);
            } else if (state === TestRunState.FAILED) {
                content.push(<Icon name="times" className="failed" key={content.length} about="Failed" />);
            } else {
                content.push(<Icon name="circle outline" className="undecided" key={content.length} about="Undecided" />);
            }
        }
        for(const c of claims) {
            if (!resolvedClaims.has(c.getName())) {
                content.push(<Icon name="user" key={content.length} about={c.getName()} />);
            }
        }
        for(let i = 0; i < remainingCellCount; i++) {
            content.push(<Icon name="user outline" key={content.length} about="Unclaimed" />);
        }
        return <div className="case-status-view">{content}</div>
    }

}