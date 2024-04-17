const router = @import("router.zig");

pub fn setup(r: *router.Router) !void {
    // Set up your custom routes here!
    try r.addRoute("/", @import("fileserver.zig").fileserver);
}
