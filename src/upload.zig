const std = @import("std");
const http = std.http;
const mem = std.mem;

// pub const Response = struct {
//     version: http.Version = .@"HTTP/1.1",
//     status: http.Status = .ok,
//     reason: ?[]const u8 = null,

//     transfer_encoding: ResponseTransfer = .none,

//     allocator: Allocator,
//     address: net.Address,
//     connection: Connection,

//     headers: http.Headers,
//     request: Request,

// method: http.Method,
// target: []const u8,
// version: http.Version,

// content_length: ?u64 = null,
// transfer_encoding: ?http.TransferEncoding = null,
// transfer_compression: ?http.ContentEncoding = null,

// headers: http.Headers,
// parser: proto.HeadersParser,
// compression: Compression = .none,

//     state: State = .first,

// POST upload data
pub fn upload(response: *http.Server.Response, _: mem.Allocator) void {
    const request: http.Server.Request = response.request;

    if (request.method != http.Method.POST) return;
}
