const http = @import("std").http;

pub fn index(response: *http.Server.Response) void {
    const body = "Index!";

    //TODO: Handle the errors lol
    response.transfer_encoding = .{ .content_length = body.len };
    response.headers.append("content-type", "text/plain") catch unreachable;
    response.do() catch unreachable;

    _ = response.writeAll(body) catch unreachable;
    response.finish() catch unreachable;
}
