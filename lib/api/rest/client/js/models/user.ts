import Users from "../controllers/users";

export default class User {
    public static fromRaw(data: any) {
        data ||= {}
        return new User(
            data['id'] || '',
            data['email'] || '',
            data['name'] || '',
            data['surname'] || '',
            data['meta'] || null
        )
    }

    constructor(
        public id: string,
        public email: string,
        public name: string,
        public surname: string,
        public meta: any
    ) {
    }
}