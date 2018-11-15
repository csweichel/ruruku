
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

export interface TestSuite {
    name: string
    tags: { [id: string]: string }
    cases: TestCase[]
}

export interface TestCase {
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
    claimedCases: string[]
}

type TestCaseResult = "pass" | "fall" | "undecided"

export interface TestCaseRun {
    caseName: string
    caseGroup: string
    start: Date
    tester: string
    result: TestCaseResult
    comment: string
}
