const http = @import("std").http;

const router = @import("../router.zig").Router;
const index = @import("index.zig").index;

// This is where you will add your routes!
pub fn setup(r: *router) !void {

    ////////DO NOT REMOVE//////////
    try r.addRoute("/", .{
        .method = http.Method.GET,
        .fileserver = true,
        .route = index,
    });
    ///////////////////////////////

    ////////Your routes////////////

    ///////////////////////////////
}
