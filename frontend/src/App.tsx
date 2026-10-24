
function App() {
  return (
    <main className="min-h-screen overflow-hidden bg-white text-slate-900">
      {/* Navigation */}
      <header className="border-b border-slate-100 bg-white">
        <nav className="mx-auto flex max-w-7xl items-center justify-between px-5 py-5 sm:px-8 lg:px-12">
          <a href="#" className="flex items-center gap-3" aria-label="CampusCore home">
            <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-indigo-600 text-lg font-extrabold text-white shadow-sm">
              CC
            </div>
            <div>
              <p className="text-lg font-extrabold tracking-tight">CampusCore</p>
              <p className="text-xs text-slate-500">Campus management, connected</p>
            </div>
          </a>

          <div className="hidden items-center gap-8 text-sm font-medium text-slate-600 md:flex">
            <a href="#about" className="transition hover:text-indigo-600">About</a>
            <a href="#features" className="transition hover:text-indigo-600">Features</a>
            <a href="#community" className="transition hover:text-indigo-600">Community</a>
          </div>

          <a
            href="#features"
            className="rounded-xl bg-indigo-600 px-4 py-3 text-sm font-semibold text-white shadow-sm transition hover:bg-indigo-700 sm:px-5"
          >
            Explore CampusCore
          </a>
        </nav>
      </header>

      {/* Hero */}
      <section className="relative">
        <div className="mx-auto grid max-w-7xl items-center gap-12 px-5 py-14 sm:px-8 sm:py-20 lg:grid-cols-2 lg:gap-16 lg:px-12 lg:py-24">
          <div className="max-w-2xl">
            <div className="mb-6 inline-flex items-center gap-2 rounded-full border border-indigo-100 bg-indigo-50 px-4 py-2 text-sm font-semibold text-indigo-700">
              <span className="h-2 w-2 rounded-full bg-indigo-600" />
              A new way to experience campus life
            </div>

            <h1 className="text-4xl font-extrabold leading-tight tracking-tight sm:text-5xl lg:text-6xl">
              Your campus.
              <span className="block text-indigo-600">Better connected.</span>
            </h1>

            <p className="mt-6 max-w-xl text-lg leading-8 text-slate-600">
              Bringing students, lecturers, and campus services together in
              one place, so academic life can work smarter and feel simpler.
            </p>

            <div className="mt-9 flex flex-col gap-3 sm:flex-row">
              <a
                href="#features"
                className="inline-flex items-center justify-center gap-2 rounded-xl bg-indigo-600 px-6 py-4 font-semibold text-white shadow-lg shadow-indigo-600/20 transition hover:-translate-y-0.5 hover:bg-indigo-700"
              >
                Discover CampusCore <span aria-hidden="true">→</span>
              </a>
              <a
                href="#about"
                className="inline-flex items-center justify-center rounded-xl border border-slate-200 bg-white px-6 py-4 font-semibold text-slate-700 transition hover:border-indigo-200 hover:bg-slate-50"
              >
                Learn more
              </a>
            </div>

            <div className="mt-10 flex flex-wrap items-center gap-x-6 gap-y-3 text-sm text-slate-500">
              <span className="flex items-center gap-2">
                <span className="text-lg text-indigo-600">✓</span>
                Designed for campus communities
              </span>
              <span className="flex items-center gap-2">
                <span className="text-lg text-indigo-600">✓</span>
                Built to bring services together
              </span>
            </div>
          </div>

          {/* Student photography */}
          <div className="relative mx-auto w-full max-w-xl lg:max-w-none">
            <div className="absolute -right-5 -top-5 h-32 w-32 rounded-full bg-indigo-100 sm:-right-7 sm:-top-7 sm:h-44 sm:w-44" />
            <div className="absolute -bottom-5 -left-5 h-28 w-28 rounded-full bg-amber-100 sm:-bottom-7 sm:-left-7 sm:h-36 sm:w-36" />

            <div className="relative overflow-hidden rounded-3xl bg-slate-100 shadow-2xl shadow-slate-900/10">
              <img
                src="/campus-students.png"
                alt="University students learning and collaborating together"
                className="h-[360px] w-full object-cover sm:h-[480px] lg:h-[540px]"
                fetchPriority="high"
              />

              <div className="absolute inset-0 bg-gradient-to-t from-slate-950/65 via-transparent to-transparent" />

              <div className="absolute bottom-0 left-0 right-0 p-6 text-white sm:p-8">
                <p className="text-sm font-semibold uppercase tracking-widest text-indigo-200">
                  One campus community
                </p>
                <h2 className="mt-2 max-w-sm text-2xl font-bold leading-snug sm:text-3xl">
                  Learn together. Grow together.
                </h2>
              </div>
            </div>

            <div className="absolute -left-3 top-8 rounded-2xl border border-white/70 bg-white/95 p-4 shadow-xl backdrop-blur sm:-left-8 sm:top-12">
              <div className="flex items-center gap-3">
                <div className="flex h-11 w-11 items-center justify-center rounded-xl bg-indigo-100 text-xl">
                  📚
                </div>
                <div>
                  <p className="font-bold">Learning first</p>
                  <p className="text-sm text-slate-500">Made for campus life</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Introduction */}
      <section id="about" className="bg-slate-50 px-5 py-16 sm:px-8 sm:py-20 lg:px-12">
        <div className="mx-auto grid max-w-7xl gap-6 md:grid-cols-2 md:items-center">
          <div>
            <p className="text-sm font-bold uppercase tracking-widest text-indigo-600">
              Welcome to CampusCore
            </p>
            <h2 className="mt-3 max-w-xl text-3xl font-bold tracking-tight sm:text-4xl">
              More connected campus experiences start here.
            </h2>
          </div>
          <p className="text-base leading-8 text-slate-600 sm:text-lg">
            CampusCore is being built to help institutions bring academic
            activities, campus services, and communication into a more
            connected experience. We are starting with the needs of African
            institutions, with a vision that can grow beyond borders.
          </p>
        </div>
      </section>

      {/* Feature preview */}
      <section id="features" className="px-5 py-16 sm:px-8 sm:py-20 lg:px-12">
        <div className="mx-auto max-w-7xl">
          <div className="mx-auto max-w-2xl text-center">
            <p className="text-sm font-bold uppercase tracking-widest text-indigo-600">
              One connected platform
            </p>
            <h2 className="mt-3 text-3xl font-bold tracking-tight sm:text-4xl">
              Built around campus life
            </h2>
            <p className="mt-4 leading-7 text-slate-600">
              A foundation for making everyday academic and administrative
              activities easier to access and manage.
            </p>
          </div>

          <div className="mt-10 grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            <article className="rounded-2xl border border-slate-200 bg-white p-6 transition hover:-translate-y-1 hover:shadow-lg">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-indigo-100 text-2xl">📖</div>
              <h3 className="mt-5 text-lg font-bold">Academic life</h3>
              <p className="mt-2 leading-7 text-slate-600">
                A place to bring academic information and learning-related
                services closer to students and lecturers.
              </p>
            </article>

            <article className="rounded-2xl border border-slate-200 bg-white p-6 transition hover:-translate-y-1 hover:shadow-lg">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-emerald-100 text-2xl">🤝</div>
              <h3 className="mt-5 text-lg font-bold">Campus community</h3>
              <p className="mt-2 leading-7 text-slate-600">
                Support better connections between students, lecturers,
                and the wider campus community.
              </p>
            </article>

            <article className="rounded-2xl border border-slate-200 bg-white p-6 transition hover:-translate-y-1 hover:shadow-lg sm:col-span-2 lg:col-span-1">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-amber-100 text-2xl">🏛️</div>
              <h3 className="mt-5 text-lg font-bold">Campus services</h3>
              <p className="mt-2 leading-7 text-slate-600">
                A foundation for organizing institutional services and
                everyday campus operations in one connected system.
              </p>
            </article>
          </div>
        </div>
      </section>

      {/* Community banner */}
      <section id="community" className="px-5 pb-16 sm:px-8 sm:pb-20 lg:px-12">
        <div className="mx-auto max-w-7xl overflow-hidden rounded-3xl bg-indigo-700 px-6 py-12 text-white sm:px-12 sm:py-16">
          <div className="max-w-3xl">
            <p className="text-sm font-bold uppercase tracking-widest text-indigo-200">
              The campus of tomorrow
            </p>
            <h2 className="mt-3 text-3xl font-bold tracking-tight sm:text-4xl">
              Every campus has a story. Let’s connect yours.
            </h2>
            <p className="mt-4 max-w-2xl leading-7 text-indigo-100">
              CampusCore is laying the groundwork for a more connected
              academic experience, starting with the people who make campus
              life matter.
            </p>
            <a
              href="#about"
              className="mt-7 inline-flex rounded-xl bg-white px-6 py-3.5 font-semibold text-indigo-700 transition hover:bg-indigo-50"
            >
              Explore our vision
            </a>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-slate-200 px-5 py-7 sm:px-8 lg:px-12">
        <div className="mx-auto flex max-w-7xl flex-col gap-3 text-sm text-slate-500 sm:flex-row sm:items-center sm:justify-between">
          <p className="font-bold text-slate-700">CampusCore</p>
          <p>Connect your campus. Empower learning.</p>
          <p>Campus management for a connected future.</p>
        </div>
      </footer>
    </main>
  )
}

export default App
