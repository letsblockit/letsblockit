# Shutting down the letsblock.it project and its official instance

The letsblock.it project enables users to remove low-quality content and useless nags, to tailor their online
experience and focus on what matters. It has been running for more than two years now and its official instance
currently serves more than 800 active users and hundreds of anonymous visitors every day. We achieved this
thanks to [dozens of contributors](https://github.com/letsblockit/letsblockit/?tab=readme-ov-file#thanks-to-our-contributors)
and [financial sponsors](https://opencollective.com/letsblockit#section-contributors) who I am grateful for.

I started this project in 2021 as a way to help people regain agency over their online experience:
[enshittification](https://pluralistic.net/2023/01/21/potemkin-ai/#hey-guys) of the commercial web had already
started, but [deceptive patterns](https://www.deceptive.design/) and user-hostile features were not so prevalent yet.
In that context, I naively hoped to be able to contain most of this through content filters and started writing the
templates that became letsblock.it.

Unfortunately, the commercial web is getting worse every day, with:

  - user-hostile design becoming the norm, as big tech corporations are set to extract as much money from society as possible
  - the "generative AI" bubble wasting precious water and electricity to drown good content into a sea of derivative low-quality content
  - the advertising company Google continuing their efforts to deny people control over their browser, neutering
    content-blocking extensions with MV3 [under false security claims](https://www.eff.org/deeplinks/2021/12/chrome-users-beware-manifest-v3-deceitful-and-threatening)
    and planning to [lock down the OS and browser with DRM](https://arstechnica.com/gadgets/2023/07/googles-web-integrity-api-sounds-like-drm-for-the-web/).

This project is making the commercial web more bearable, but I'd rather spend my energy on making the non-commercial
web more attractive. I want to support communities and applications that respect their users and value what we
have to say. These websites don't need letsblock.it rules, because they don't shove low-quality content and anti-features down our throats.


## What will happen

- The official instance hosted at https://letsblock.it/ will be shut down. Because its database contains
  personal information, I will not hand over its operations to a third party. Its database will be
  destroyed after users get a chance to migrate out.
- Limited maintenance to guarantee security will continue for a few months, but the template corpus will not
  receive and fix or update. 
- If a group of users wants to carry these efforts forward, I'll support them in forking the project and setting
  up a new server if they wish. I'll happily publicize this fork and assist users in migrating their data.

## Provisional timeline

#### March 2024: official announcement
  - Announcement on OpenCollective and removal of all contribution options
  - Announcement on a GitHub discussion

### April 2024: 
  - Registrations of new accounts on the official instance will be disabled, notice is added on top of all pages
  - Import and export functions will be added and documented
  - Recommended alternatives will be listed to help users move on

#### June 2024: shutdown of the official instance, archival of the project
  - The official instance will be shutdown and all user data deleted at the end of the month
  - The GitHub project will be archived and not further updates released, even for security fixes

## What's next?

After several failed projects and abandoned hobbies, launching letsblock.it and keeping it running for over
two years is a big success in my book. I am very grateful to everyone who contributed to this project and
made it possible. My motivation has always been to be useful to others and I think it's better to give our
users a controlled shutdown than keep the project stale and unmaintained.

This is not a call for contributors or new maintainers. If users want to take over this project's legacy,
I'd prefer it to be done under a new name to avoid any confusion. Feel free to coordinate on this discussion
though, and I'd be happy to support you with the setup.



## What's next for users

- Existing users can go to their [user preferences page](/user/account) to download a static version
  of their rules, or export their list to use it with another instance.

- https://iorate.github.io/ublacklist/  https://github.com/quenhus/uBlock-Origin-dev-filter
- https://addons.mozilla.org/fr/firefox/search/?q=youtube
- https://www.tampermonkey.net/
- 
