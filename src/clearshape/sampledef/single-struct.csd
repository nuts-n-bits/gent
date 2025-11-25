type something_something = {
    s session: string,
    t timestamp: {
        i isoTimestamp: i64,
        u unixTimestampMillis: u64,
    },
    u user: {
        name: string,
        friends: string,
    }
}


type ArSync = {
    s session: {
        n caseNumber: string,
        i caseId: string,
    },
    r revisions: {
        r revid: string,
        t timestamp: {
            i isoTimestamp: string,
            u unixTimestampMillis: u64,
        },
        u userid: string,
        n username: string,
        s editSummary: string,
        c content: {
            s string: string,
        }
    }
}