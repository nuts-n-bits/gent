type something_something = {
    s session?: string,
    t timestamp: enum {
        i isoTimestamp: [string, i64],
        u unixTimestampMillis: u64,
    },
    u user: {
        name: string,
        friends: {
            id: string,
        }[]
    }
}


type ArSync = {
    s session: {
        n caseNumber: string,
        i caseId: string,
    },
    r revisions: {
        r revid: string,
        t timestamp: enum {
            i isoTimestamp: [string, i64],
            u unixTimestampMillis: u64,
        }
        u userid: string,
        n username: string,
        s editSummary: string[],
        c content: enum {
            s string: string,
            l segmentList: map(string),
        }
    }
}