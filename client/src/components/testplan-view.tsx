import * as React from 'react';
import { Table, Button, ButtonGroup, ButtonOr, Checkbox, Menu, Modal, Icon } from 'semantic-ui-react';
import { TestRunStatus, TestcaseStatus, TestcaseRunResult, TestRunState, Testcase } from '../api/v1/session_pb';
import { TestcaseDetailView } from './testcase-detail-view';
import { TestCaseStatusView } from './testcase-status'
import { NewTestcaseRunView } from './new-testcase-run-view';
import { AppStateContent, getClient } from 'src/types/app-state';
import { Disposable } from 'src/types/service-client';
import { NewTestcaseView } from './new-testcase-view';

import './testplan-view.css';

export interface TestplanViewProps {
    appState: AppStateContent
    showDetails(content?: any): void
}

type SortableColumns = "group" | "name" | "actions";
interface TestplanViewState {
    column: SortableColumns
    direction: 'ascending' | 'descending'
    status?: TestRunStatus
    editMode: boolean
    testcaseToRemove?: Testcase
}

export class TestplanView extends React.Component<TestplanViewProps, TestplanViewState> {
    protected updateListener: Disposable;
    protected lastOpenedTestcaseIndex: number | undefined;

    public constructor(props: TestplanViewProps) {
        super(props);
        this.state = {
            column: 'name',
            direction: 'ascending',
            editMode: false
        };

        this.closeSidebar = this.closeSidebar.bind(this);
        this.toggleEditMode = this.toggleEditMode.bind(this);
    }

    public async componentWillMount() {
        this.props.appState.keybindings.push(
            { keys: 'esc', description: 'Close the sidebar', handler: () => this.closeSidebar() },
            { keys: 'e', description: "Toggle edit mode", handler: this.toggleEditMode },
            { keys: 'b', description: "Mark current testcase as passed", handler: this.claimContributeAdvance.bind(this, TestRunState.PASSED) },
            { keys: 'n', description: "Mark current testcase as undecided", handler: this.claimContributeAdvance.bind(this, TestRunState.UNDECIDED) },
            { keys: 'm', description: "Mark current testcase as failed", handler: this.claimContributeAdvance.bind(this, TestRunState.FAILED) }
        );

        this.fetchStatus();
        try {
            this.updateListener = await getClient(this.props.appState).listenForUpdates((s, err) => {
                if (err) {
                    this.props.appState.setError(`${err}`);
                } else {
                    this.setState({ status: s });
                }
            });
        } catch (err) {
            this.props.appState.setError(err);
        }
    }

    public async componentWillUnmount() {
        this.props.appState.keybindings.pop();
        if (this.state.editMode) {
            this.props.appState.keybindings.pop();
        }
    }

    public render() {
        const { column, direction } = this.state;

        const testcases = this.buildRows();
        return (
            <div>
            <Modal open={!!this.state.testcaseToRemove}>
                <Modal.Header>Remove testcase</Modal.Header>
                <Modal.Content>
                    <p>You are about to remove the testcase {this.state.testcaseToRemove ? this.state.testcaseToRemove.getId() : ""}.
                    <b>This action cannot be undone.</b> Do you want to continue?</p>
                </Modal.Content>
                <Modal.Actions>
                    <Button basic={true} color='red' onClick={this.removeTestcase.bind(this, undefined)}>Cancel</Button>
                    <Button color='green' inverted={true} onClick={this.removeTestcase.bind(this, this.state.testcaseToRemove)}><Icon name='checkmark' /> Do it</Button>
                </Modal.Actions>
            </Modal>
            {this.state.status && this.state.status.getModifiable() && (
                <Menu secondary={true}>
                    {this.state.editMode && (
                        <Menu.Item onClick={this.showNewTestcaseForm.bind(this, undefined)}><Icon name="plus" />Add testcase</Menu.Item>
                    )}
                    <Menu.Menu position="right">
                        <Menu.Item><Checkbox toggle={true} label="Edit test plan" checked={this.state.editMode} onChange={this.toggleEditMode} /></Menu.Item>
                    </Menu.Menu>
                </Menu>
            )}
            <Table celled={true} sortable={true} fixed={true}>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell sorted={column === 'group' ? direction : undefined} onClick={this.handleSort.bind(this, 'group')}>Group</Table.HeaderCell>
                        <Table.HeaderCell sorted={column === 'name' ? direction : undefined} onClick={this.handleSort.bind(this, 'name')}>Name</Table.HeaderCell>
                        <Table.HeaderCell>Testers</Table.HeaderCell>
                        <Table.HeaderCell sorted={column === 'actions' ? direction : undefined} onClick={this.handleSort.bind(this, 'actions')} />
                    </Table.Row>
                </Table.Header>
                <Table.Body>
                    {testcases}
                </Table.Body>
            </Table>
            </div>
        )
    }

    protected async fetchStatus() {
        try {
            const s = await getClient(this.props.appState).getStatus();
            this.setState({ status: s });
            if (!!s && s.getModifiable() && s.getCaseList().length === 0) {
                this.toggleEditMode();
            }
        } catch (err) {
            this.props.appState.setError(err);
        }
    }

    protected handleSort(column: SortableColumns) {
        if (this.state.column === column) {
            this.setState({direction: this.state.direction === 'ascending' ? 'descending' : 'ascending'});
        } else {
            this.setState({
                column,
                direction: 'ascending'
            })
        }
    }

    protected buildRows() {
        if (!this.state.status) {
            return [];
        }

        const cases = this.state.status.getCaseList();
        if (!cases) {
            console.warn("Session has empty case list");
            return [];
        }

        const { column, direction } = this.state;
        const sortedCases = cases.sort((a, b) => {
            let result = 0;

            if (column === 'actions') {
                if (this.isClaimed(a)) {
                    result = 1;
                } else if (this.isClaimed(b)) {
                    result = -1;
                } else {
                    result = 0;
                }
            } else if (column === 'group') {
                result =  a.getCase()!.getGroup()!.localeCompare(b.getCase()!.getGroup());
            } else if (column === 'name') {
                result = a.getCase()!.getName()!.localeCompare(b.getCase()!.getName());
            }
            if (direction === 'descending') {
                result *= -1;
            }
            return result;
        });

        return sortedCases.map((tcs, idx) => {
            const cse = tcs.getCase()!;
            const caseComplete = tcs.getResultList().length >= cse.getMintestercount();
            const positive = caseComplete && tcs.getState() === TestRunState.PASSED;
            const warning = caseComplete && tcs.getState() === TestRunState.UNDECIDED;
            const negative = caseComplete && tcs.getState() === TestRunState.FAILED;
            return <Table.Row key={cse.getId()} active={this.lastOpenedTestcaseIndex === idx} positive={positive} warning={warning} negative={negative}>
                <Table.Cell>{cse.getGroup()}</Table.Cell>
                <Table.Cell><a href="#" onClick={this.showDetails.bind(this, idx, tcs)}>{cse.getName()}</a></Table.Cell>
                <Table.Cell><TestCaseStatusView case={tcs} /></Table.Cell>
                <Table.Cell collapsing={true}>
                    {this.getActions(idx, tcs)}
                </Table.Cell>
            </Table.Row>
        });
    }

    protected getActions(index: number, tc: TestcaseStatus) {
        if (this.state.editMode) {
            return (
                <ButtonGroup>
                    <Button icon="pencil" key="edit" onClick={this.showNewTestcaseForm.bind(this, tc.getCase())} />
                    <Button icon="minus" key="remove" onClick={this.removeTestcase.bind(this, tc.getCase())} />
                </ButtonGroup>
            );
        }

        if (this.isClaimed(tc)) {
            const previousRun = this.getPreviousRun(tc);
            if (previousRun) {
                return <Button label="Edit contribution" icon="write square" key="contribute" onClick={this.showNewRunForm.bind(this, index, tc, undefined)} />;
            } else {
                return (
                    <ButtonGroup>
                        <Button icon="check" key="pass" onClick={this.showNewRunForm.bind(this, index, tc, TestRunState.PASSED)} />
                        <Button icon="question" key="undecided" onClick={this.showNewRunForm.bind(this, index, tc, TestRunState.UNDECIDED)} />
                        <Button icon="times" key="fail" onClick={this.showNewRunForm.bind(this, index, tc, TestRunState.FAILED)} />
                        <ButtonOr />
                        <Button label="Unclaim" icon="minus circle" onClick={this.claim.bind(this, tc, false)} key="claim" />;
                    </ButtonGroup>
                );
            }
        } else {
            return <Button label="Claim" icon="plus circle" onClick={this.claim.bind(this, tc, true)} key="claim" />;
        }
    }

    protected showDetails(idx: number, cse: TestcaseStatus, evt: React.SyntheticEvent) {
        evt.preventDefault();
        this.lastOpenedTestcaseIndex = idx;
        this.props.showDetails(<TestcaseDetailView testcase={cse} onClose={this.closeSidebar} />);
    }

    protected showNewRunForm(idx: number, cse: TestcaseStatus, result?: TestRunState) {
        this.lastOpenedTestcaseIndex = idx;
        this.props.showDetails(<NewTestcaseRunView
            appState={this.props.appState}
            testcase={cse}
            result={result}
            onClose={this.closeSidebar} />)
    }

    protected showNewTestcaseForm(cse?: Testcase) {
        this.props.showDetails(<NewTestcaseView
            appState={this.props.appState}
            testcase={cse}
            onClose={this.closeSidebar} />)
    }

    protected isClaimed(cse: TestcaseStatus): boolean {
        return !!cse.getClaimList().find(c => c.getName() === this.props.appState.user!.name);
    }

    protected getPreviousRun(cse: TestcaseStatus): TestcaseRunResult | undefined {
        return cse.getResultList().find(r => {
            const p = r.getParticipant();
            return !!(p && p.getName() === this.props.appState.user!.name);
        });
    }

    protected async claim(cse: TestcaseStatus, claim: boolean, evt: React.SyntheticEvent): Promise<void> {
        evt.preventDefault();
        try {
            await getClient(this.props.appState).claim(cse.getCase()!.getId(), claim);
        } catch (err) {
            this.props.appState.setError(err);
        }
    }

    protected closeSidebar() {
        this.lastOpenedTestcaseIndex = undefined;
        this.props.showDetails(undefined);
    }

    protected async claimContributeAdvance(result: TestRunState): Promise<void> {
        if (!this.state.status) {
            return;
        }

        const cases = this.state.status.getCaseList();
        if (!cases || this.lastOpenedTestcaseIndex === undefined || this.lastOpenedTestcaseIndex >= cases.length) {
            return;
        }

        const cse = cases[this.lastOpenedTestcaseIndex];
        try {
            const client = await getClient(this.props.appState);
            await client.claim(cse.getCase()!.getId(), true);
            await client.contribute(cse.getCase()!.getId(), result, "");

            const nidx = this.lastOpenedTestcaseIndex + 1;
            if (nidx >= cases.length) {
                this.closeSidebar();
            } else {
                this.showNewRunForm(nidx, cases[nidx], result);
            }
        } catch (err) {
            this.props.appState.setError(err);
        }
    }

    protected toggleEditMode() {
        let newEditMode = !this.state.editMode;
        if (!this.state.status || !this.state.status.getModifiable()) {
            newEditMode = false;
        }

        this.setState({ editMode: newEditMode });
        this.lastOpenedTestcaseIndex = undefined;

        if (newEditMode) {
            this.props.appState.keybindings.push({ keys: 'a', description: "Add a new testcase", handler: () => this.showNewTestcaseForm(undefined) });
        } else {
            this.props.appState.keybindings.pop();
        }
    }

    protected async removeTestcase(tc?: Testcase) {
        if (tc && tc === this.state.testcaseToRemove) {
            this.setState({ testcaseToRemove: undefined });
            try {
                const client = await getClient(this.props.appState);
                await client.removeTestcase(tc);
            } catch (err) {
                this.props.appState.setError(err);
            }
        } else if (tc) {
            // open the modal
            this.setState({ testcaseToRemove: tc });
        } else {
            // close the modal
            this.setState({ testcaseToRemove: undefined });
        }
    }

}