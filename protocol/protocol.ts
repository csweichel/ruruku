
export interface WelcomeRequest {
    type: "welcome"
    name: string
}

export interface WelcomeResponse {
    type: "welcome"
    suite: TestSuite
    run: TestRun
    participant: TestParticipant
}

export interface KeepAliveRequest {
    type: "keep-alive"
}

export interface ClaimRequest {
    type: "claim"
    caseId: string
    claim: boolean
}

export interface ClaimResponse {
    type: "claim"
}

export interface UpdateMessage {
    type: "update"
    run: TestRun
    participant: TestParticipant
}


export interface TestSuite {
    name: string
    tags: { [id: string]: string }
    cases: TestCase[]
}

export interface TestCase {
    id: string
    name: string
    group: string
    description: string
    steps: string
    mustPass: boolean
    minTesterCount: number
}


export interface TestRun {
    suiteName: string
    start: Date
    participants: TestParticipant[]
    cases: TestCaseRun[]
}

export interface TestParticipant {
    name: string
    claimedCases: { [id: string]: boolean }
}

type TestCaseResult = "pass" | "fall" | "undecided"

export interface TestCaseRun {
    case: string
    caseGroup: string
    start: Date
    tester: string
    result: TestCaseResult
    comment: string
}
