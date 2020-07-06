# Hegemonie

[![CircleCI](https://circleci.com/gh/jfsmig/hegemonie.svg?style=svg)](https://circleci.com/gh/jfsmig/hegemonie)
[![Codecov](https://codecov.io/gh/jfsmig/hegemonie/branch/master/graph/badge.svg)](https://codecov.io/gh/jfsmig/hegemonie)
[![Codacy](https://app.codacy.com/project/badge/Grade/bf7c2872c60445c99f914d31d7b213ae)](https://www.codacy.com/manual/jfsmig/hegemonie?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=jfsmig/hegemonie&amp;utm_campaign=Badge_Grade)
[![MPL-2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

Hegemonie is an online management, strategy and diplomacy RPG. The current
repository is a reboot of what [Hegemonie](http://www.hegemonie.be) was between
1999 and 2003. It is under "heavily inactive" construction.

Learn more details on the [GAME description](./GAME.md) page.

## Getting Started

Start a sandbox environment using ``docker-compose``. For more information,
please refer to the page with the [technical elements](./TECH.md).

```shell
BASE=github.com/jfsmig/hegemonie
go get "${BASE}"
cd "${GOPATH}/${BASE}"
docker-compose up
```

## How To Contribute

Contributions are what make the open source community such an amazing place. Any
contributions you make are greatly appreciated.

1. Fork the Project
2. Create your Feature Branch (git checkout -b feature/MyFeature)
3. Commit your Changes (git commit -m 'Add some MyFeature')
4. Push to the Branch (git push origin feature/MyFeature)
5. Open a Pull Request... and if it takes to long to get a feedback, that's
   probably because the maintainers didn't notice it. Do not hesitate to pig
   them :).

## License

Hegemonie is distributed under the MPLv2 License. See [LICENSE](./LICENSE) for
more information.

We strongly believe in Open Source for many reasons:

* For the purpose of a better user experience, because bugs will happen and
  opening the code is the best way to let a skilled user find it as soon as
  possible. The value of a game instance definitely is in its players and in the
  description of the world.
* For software quality purposes because a software with open sources is the best
  way to have its bugs identified and fixed as soon as possible.
* For a greater adoption, we deliberately chose a liberal license so that there
  cannot be any legal concern in using the code.
* For transparency reasons. an open code lets you check that nothing odd is done
  with your data.

## Contact

Follow the development on GitHub with the
[jfsmig/hegemonie](https://github.com/jfsmig/hegemonie) project. A Facebook page
also exists [Hegemonie.be](https://www.facebook.com/hegemonie.be).

## Acknowledgements

We welcome any volunteer and we already have a list of
[amazing authors of Hegemonie](./AUTHORS.md).
