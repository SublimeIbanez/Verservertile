
export class Result {
    constructor(okVal = null, errVal = null) {
        this.Ok = okVal;
        this.Err = errVal;
    }

    static Ok = (value) => {
        return new Result(value, null);
    }

    static Err = (value) => {
        return new Result(null, value);
    }

    isOk() {
        return this.Ok !== null;
    }

    isErr() {
        return this.Err !== null;
    }

    match({ Ok, Err }) {
        return this.isOk() ? Ok(this.Ok) : Err(this.Err);
    }

    unwrap() {
        return this.Ok;
    }
};