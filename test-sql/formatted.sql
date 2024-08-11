create table foo_bar_baz (
    id       serial         primary key,
    foo      varchar(255)   not null,
    bar      integer        not null,
    baz      date,
    qux      boolean        default false,
    quux     uuid           not null default 123 foo_foo(id),
    holymoly vectors.vector not null
);

select
    foo,
    bar,
    baz
from
    foo_bar_baz
where
    qux = true
order by
    baz desc;
