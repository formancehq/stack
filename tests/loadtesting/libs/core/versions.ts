const semverRegexp = /^v?(?<Major>0|[1-9]\d*)\.(?<Minor>0|[1-9]\d*)\.(?<Patch>0|[1-9]\d*)(?:-(?<PreRelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?<Meta>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$/;

export class Version {
    constructor(
        public readonly major: number,
        public readonly minor: number,
        public readonly patch: number,
    ) {
    }

    lt(v: Version) {
        return this.major < v.major ||
            this.minor < v.minor ||
            this.patch < v.patch;
    }
}

export const parseVersion = (v: string) => {
    if (!semverRegexp.test(v)) {
        return new Version(1, 7, 4);
        // throw new Error(`Invalid version ${v}, not matching semver`);
    }
    const ret = semverRegexp.exec(v);
    const major = ret!.groups!.Major as string;
    const minor = ret!.groups!.Minor;
    const patch = ret!.groups!.Patch;

    return new Version(
        parseInt(major, 10),
        minor ? parseInt(minor, 10) : 0,
        patch ? parseInt(patch, 10) : 0
    );
};
