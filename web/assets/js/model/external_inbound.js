class DBExternalInbound {
    constructor(data) {
        this.id = 0;
        this.enable = true;
        this.name = "";
        this.comment = "";
        this.url = "";
        this.clients = [];
        if (data == null) {
            return;
        }
        ObjectUtil.cloneProps(this, data);
    }
}
