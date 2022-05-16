#!/opt/homebrew/bin/python3

import sys
import json

def main():
    jsonStr = sys.stdin.readline()
    aDict = json.loads(jsonStr)

    response = {"version": "v1", "isWhiteOut": False, "patches":[]}

    if len(aDict) == 0:
        metadata = {"name": "PythonTesting", "version":"v0.0.1","requestVersion":["v1"],"responseVersion":["v1"]}
        print(json.dumps(metadata, indent=4))
        return

    patches = []
    # getting the containers
    try:
        containers = aDict['spec']['template']['spec']['containers']
        for i, c in enumerate(containers):
            if "securityContext" in c:
                patch = {"op": "remove", "path": f'/spec/template/spec/containers/{i}/securityContext'}
                patches.append(patch)
    except KeyError:
        pass

    try:
        securityContext = aDict['spec']['template']['spec']['securityContext']
        patches.append({"op": "remove", "path": "/spec/template/spec/securityContext"})
    except KeyError:
        pass

    response['patches'] = patches
    print(json.dumps(response, indent=4))



if __name__ == "__main__":
    main()
