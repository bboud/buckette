const http = @import("std").http;

const router = @import("../router.zig").Router;
const index = @import("../fileserver.zig").fileserver;

// This is where you will add your routes!
pub fn setup(r: *router) !void {

    ////////DO NOT REMOVE//////////
    try r.addRoute("/", index);
    ///////////////////////////////

    ////////Your routes////////////

    ///////////////////////////////
}
